# -*- mode: ruby -*-

Vagrant.configure(2) do |config|
  dirs = ENV["GOPATH"] || Dir.home

  config.vm.box = "boxcutter/ubuntu1504-docker"
  config.vm.network "forwarded_port", guest: 2375, host: 12375
  config.vm.host_name = "devbox"
  config.vm.synced_folder ".", "/vagrant", disabled: true

  dirs.split(File::PATH_SEPARATOR).each do |dir|
    config.vm.synced_folder dir, dir
  end

  [:vmware_fusion, :vmware_workstation].each do |visor|
    config.vm.provider visor do |v, override|
      v.vmx["memsize"] = "2048"
    end
  end

  Dir["machines/devbox/provision*.sh"].each do |path|
    config.vm.provision "shell", path: path
  end
end
