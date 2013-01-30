require "rubygems"
require "json"

def abs_path(path)
  File.expand_path(File.join(current_path, path))
end

def current_path
  File.expand_path(File.dirname(__FILE__))
end

def go_command(command)
  sh "env GOPATH=#{current_path} #{command}"
end

def dependencies
  @dependencies ||= JSON.parse(File.read(abs_path("dependencies.json")))
end

desc "clean compiled files and binaries"
task :clean do
  rm_rf abs_path("bin")
  rm_rf abs_path("pkg")
end

desc "fetch dependencies"
task :deps do
  deps_to_fetch = dependencies.reject do |dependency|
    path = abs_path("src/#{dependency}")
    File.exist?(path)
  end

  deps_to_fetch.each do |dependency|
    go_command("go get #{dependency}")
    rm_rf abs_path("src/#{dependency}/.git")
  end
end

namespace :deps do
  desc "Update a package"
  task :update, :package do |t, args|
    rm_rf abs_path("src/#{args[:package]}/")
    go_command("go get #{args[:package]}")
    rm_rf abs_path("src/#{args[:package]}/.git")
  end
end

desc "format code"
task :fmt do
  go_command("go fmt netracker/...")
end

desc "build the project"
task :build => [:clean, :deps, :fmt] do
  dependencies.each { |dep| go_command("go install #{dep}") }
  go_command("go install netracker/...")
end

desc "run the tests"
task :test => [:build] do
  go_command("go test netracker/...")
end

task :default => :test
