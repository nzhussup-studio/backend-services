#!/usr/bin/env ruby

require 'json'
require 'yaml'

def deep_transform(value, &block)
  case value
  when Hash
    value.each_with_object({}) { |(key, child), memo| memo[key] = deep_transform(child, &block) }
  when Array
    value.map { |item| deep_transform(item, &block) }
  else
    block.call(value)
  end
end

def prefix_map(hash, prefix)
  return {} unless hash.is_a?(Hash)

  hash.each_with_object({}) do |(key, value), memo|
    memo["#{prefix}_#{key}"] = value
  end
end

def adjust_security_requirements!(value, security_names)
  case value
  when Hash
    if value.key?('security') && value['security'].is_a?(Array)
      value['security'] = value['security'].map do |item|
        next item unless item.is_a?(Hash)

        item.each_with_object({}) do |(key, child), memo|
          memo[security_names.fetch(key, key)] = child
        end
      end
    end

    value.each_value { |child| adjust_security_requirements!(child, security_names) }
  when Array
    value.each { |child| adjust_security_requirements!(child, security_names) }
  end
end

def rewrite_refs(value, schema_names, security_names)
  deep_transform(value) do |leaf|
    next leaf unless leaf.is_a?(String)

    rewritten = leaf.gsub(%r{#/definitions/([^/]+)}) do
      name = Regexp.last_match(1)
      "#/components/schemas/#{schema_names.fetch(name, name)}"
    end

    rewritten.gsub(%r{#/components/schemas/([^/]+)}) do
      name = Regexp.last_match(1)
      "#/components/schemas/#{schema_names.fetch(name, name)}"
    end
  end.tap do |transformed|
    adjust_security_requirements!(transformed, security_names)
  end
end

def normalize_doc(service_name, doc)
  prefix = service_name.gsub(/[^A-Za-z0-9]+/, '_')
  schemas = doc.dig('components', 'schemas') || doc['definitions'] || {}
  security_schemes = doc.dig('components', 'securitySchemes') || doc['securityDefinitions'] || {}

  schema_names = schemas.keys.each_with_object({}) { |name, memo| memo[name] = "#{prefix}_#{name}" }
  security_names = security_schemes.keys.each_with_object({}) { |name, memo| memo[name] = "#{prefix}_#{name}" }

  normalized = rewrite_refs(doc, schema_names, security_names)
  paths = normalized['paths'] || {}
  components = normalized['components'] || {}

  merged_components = {
    'schemas' => prefix_map(components['schemas'] || normalized['definitions'], prefix),
    'securitySchemes' => prefix_map(components['securitySchemes'] || normalized['securityDefinitions'], prefix),
    'responses' => prefix_map(components['responses'], prefix),
    'parameters' => prefix_map(components['parameters'], prefix),
    'requestBodies' => prefix_map(components['requestBodies'], prefix),
    'headers' => prefix_map(components['headers'], prefix)
  }.delete_if { |_key, value| value.nil? || value.empty? }

  {
    'paths' => paths,
    'components' => merged_components,
    'tags' => normalized['tags'] || []
  }
end

def merge_named_hash!(target, source, section)
  source.each do |name, value|
    if target.key?(name) && target[name] != value
      abort("Error: conflicting #{section} entry #{name} while building unified OpenAPI")
    end

    target[name] = value
  end
end

manifest_file = ARGV[0]
repo_root = ARGV[1]

abort('Error: missing manifest file argument') if manifest_file.nil? || manifest_file.empty?
abort('Error: missing repo root argument') if repo_root.nil? || repo_root.empty?
abort("Error: manifest not found: #{manifest_file}") unless File.exist?(manifest_file)

manifest = JSON.parse(File.read(manifest_file))
services = manifest['services']

unless services.is_a?(Array) && services.all? { |item| item.is_a?(String) && !item.strip.empty? }
  abort('Error: openapi-services.json must contain a non-empty string array in services')
end

unified = {
  'openapi' => '3.0.1',
  'info' => {
    'title' => 'Backend Services API',
    'version' => '1.0.0'
  },
  'paths' => {},
  'components' => {},
  'tags' => []
}

seen_tags = {}

services.each do |service|
  spec_path = File.join(repo_root, service, 'openapi.yaml')
  abort("Error: generated OpenAPI file not found: #{spec_path}") unless File.exist?(spec_path)

  parsed = YAML.safe_load(File.read(spec_path), aliases: true)
  abort("Error: failed to parse #{spec_path}") unless parsed.is_a?(Hash)

  normalized = normalize_doc(service, parsed)

  normalized['paths'].each do |path, path_item|
    if unified['paths'].key?(path)
      existing = unified['paths'][path]
      path_item.each do |method, operation|
        if existing.key?(method) && existing[method] != operation
          abort("Error: duplicate operation #{method.upcase} #{path} while building unified OpenAPI")
        end

        existing[method] = operation
      end
    else
      unified['paths'][path] = path_item
    end
  end

  normalized['components'].each do |section, entries|
    unified['components'][section] ||= {}
    merge_named_hash!(unified['components'][section], entries, section)
  end

  normalized['tags'].each do |tag|
    next unless tag.is_a?(Hash) && tag['name']
    next if seen_tags[tag['name']]

    seen_tags[tag['name']] = true
    unified['tags'] << tag
  end
end

File.write(File.join(repo_root, 'openapi.yaml'), YAML.dump(unified))
