#!/usr/bin/env ruby

require 'json'

manifest_file = ARGV[0]
abort('Error: missing manifest file argument') if manifest_file.nil? || manifest_file.empty?
abort("Error: manifest not found: #{manifest_file}") unless File.exist?(manifest_file)

manifest = JSON.parse(File.read(manifest_file))
services = manifest['services']

unless services.is_a?(Array) && services.all? { |item| item.is_a?(String) && !item.strip.empty? }
  abort('Error: openapi-services.json must contain a non-empty string array in services')
end

services.each { |service| puts service }
