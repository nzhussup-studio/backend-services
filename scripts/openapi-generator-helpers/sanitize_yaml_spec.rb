#!/usr/bin/env ruby

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

def rewrite_swagger_refs(value)
  deep_transform(value) do |leaf|
    next leaf unless leaf.is_a?(String)
    leaf.gsub(%r{#/definitions/([^/]+)}) { "#/components/schemas/#{Regexp.last_match(1)}" }
  end
end

def normalize_swagger_schema(schema)
  normalized = rewrite_swagger_refs(schema || {})
  return normalized unless normalized.is_a?(Hash)

  if normalized['type'] == 'file'
    normalized = normalized.dup
    normalized['type'] = 'string'
    normalized['format'] = 'binary'
  end

  normalized
end

def convert_form_data_params(params)
  schema = { 'type' => 'object', 'properties' => {} }
  required = []

  params.each do |param|
    property =
      if param['type'] == 'file'
        { 'type' => 'string', 'format' => 'binary' }
      else
        {
          'type' => param['type'] || 'string'
        }
      end

    property['description'] = param['description'] if param['description']
    property['enum'] = param['enum'] if param['enum']
    property['items'] = normalize_swagger_schema(param['items']) if param['items']
    schema['properties'][param['name']] = property
    required << param['name'] if param['required']
  end

  schema['required'] = required unless required.empty?
  schema
end

def convert_parameter(param)
  converted = rewrite_swagger_refs(param.dup)
  schema = {}

  %w[type format enum default items].each do |key|
    next unless converted.key?(key)
    schema[key] = converted.delete(key)
  end

  schema = normalize_swagger_schema(schema) unless schema.empty?

  converted['schema'] = schema unless schema.empty?
  converted
end

def convert_response(response, produces)
  converted = response.dup
  schema = converted.delete('schema')
  return rewrite_swagger_refs(converted) unless schema

  media_types = Array(produces)
  media_types = ['application/json'] if media_types.empty?
  converted['content'] = media_types.each_with_object({}) do |media_type, memo|
    memo[media_type] = { 'schema' => normalize_swagger_schema(schema) }
  end
  rewrite_swagger_refs(converted)
end

def convert_operation(operation, default_consumes, default_produces)
  converted = rewrite_swagger_refs(operation)
  consumes = Array(converted.delete('consumes'))
  produces = Array(converted.delete('produces'))
  consumes = default_consumes if consumes.empty?
  produces = default_produces if produces.empty?

  parameters = Array(converted.delete('parameters'))
  body_param = parameters.find { |param| param['in'] == 'body' }
  form_data_params = parameters.select { |param| param['in'] == 'formData' }
  remaining_params = parameters.reject { |param| param['in'] == 'body' || param['in'] == 'formData' }

  unless remaining_params.empty?
    converted['parameters'] = remaining_params.map { |param| convert_parameter(param) }
  end

  request_body = {}
  if body_param
    media_type = consumes.first || 'application/json'
    request_body['required'] = true if body_param['required']
    request_body['description'] = body_param['description'] if body_param['description']
    request_body['content'] = {
      media_type => {
        'schema' => normalize_swagger_schema(body_param['schema'] || {})
      }
    }
  end

  unless form_data_params.empty?
    request_body['required'] = true if form_data_params.any? { |param| param['required'] }
    request_body['content'] ||= {}
    request_body['content']['multipart/form-data'] = {
      'schema' => convert_form_data_params(form_data_params)
    }
  end

  converted['requestBody'] = request_body unless request_body.empty?

  if converted['responses'].is_a?(Hash)
    converted['responses'] = converted['responses'].each_with_object({}) do |(status, response), memo|
      memo[status] = convert_response(response, produces)
    end
  end

  converted
end

def convert_swagger2_to_openapi3(doc)
  default_consumes = Array(doc.delete('consumes'))
  default_produces = Array(doc.delete('produces'))

  converted = rewrite_swagger_refs(doc)
  converted.delete('swagger')
  converted.delete('host')
  converted.delete('schemes')
  converted.delete('basePath')
  converted.delete('definitions')
  converted.delete('securityDefinitions')

  converted['openapi'] = '3.0.1'
  converted['components'] ||= {}
  converted['components']['schemas'] = normalize_swagger_schema(doc['definitions'] || {}) unless (doc['definitions'] || {}).empty?
  converted['components']['securitySchemes'] = rewrite_swagger_refs(doc['securityDefinitions'] || {}) unless (doc['securityDefinitions'] || {}).empty?

  converted['paths'] = Array(doc['paths']).empty? ? {} : doc['paths'].each_with_object({}) do |(path, path_item), memo|
    memo[path] = path_item.each_with_object({}) do |(method, operation), path_memo|
      path_memo[method] = convert_operation(operation, default_consumes, default_produces)
    end
  end

  converted
end

def reorder_top_level(doc)
  preferred = %w[openapi info paths components tags]
  reordered = {}

  preferred.each do |key|
    reordered[key] = doc[key] if doc.key?(key)
  end

  doc.each do |key, value|
    reordered[key] = value unless reordered.key?(key)
  end

  reordered
end

spec_file = ARGV[0]
abort('Error: missing spec file argument') if spec_file.nil? || spec_file.empty?
abort("Error: spec file not found: #{spec_file}") unless File.exist?(spec_file)

doc = YAML.safe_load(File.read(spec_file), aliases: true)
abort("Error: failed to parse #{spec_file}") unless doc.is_a?(Hash)

doc = convert_swagger2_to_openapi3(doc) if doc['swagger'] == '2.0' || doc['swagger'] == 2.0
doc.delete('servers')
doc.delete('host')
doc.delete('schemes')
doc = reorder_top_level(doc)

File.write(spec_file, YAML.dump(doc))
