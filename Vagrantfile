# -*- mode: ruby -*-

Vagrant.configure(2) do |config|
  dirs = ENV['GOPATH'] || Dir.home
  gdir = nil

  config.vm.box = 'boxcutter/ubuntu1504-docker'
  config.vm.network 'forwarded_port', guest: 2375, host: 12375
  config.vm.host_name = 'devbox'
  config.vm.synced_folder '.', '/vagrant', disabled: true
  config.ssh.username = 'vagrant'

  dirs.split(File::PATH_SEPARATOR).each do |dir|
    gdir = dir.sub("C\:", "/C")
    config.vm.synced_folder dir, gdir
  end

  config.vm.provider :virtualbox do |v, _override|
    v.memory = 2048
  end

  [:vmware_fusion, :vmware_workstation].each do |visor|
    config.vm.provider visor do |v, _override|
      v.memory = 2048
    end
  end

  Dir['machines/devbox/provision*.sh'].each do |path|
    config.vm.provision 'shell', path: path, args: [gdir, config.ssh.username]
  end
end
