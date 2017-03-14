Vagrant.require_version '>= 1.6.2'

# This method escapes strings for consumption inside a VMX file, used to
# inject guestinfo.ovfEnv
def vmxEscape (s)
  escaped_string = String.new
  s.each_char do |c|
    escaped_string << escape(c)
  end 
  escaped_string
end

def escape (c)
  escaped = ['#', '|', '\\', '"']
  if escaped.include? c or c.ord < 32
    return "|%02x" % c.ord
  else
    return c
  end
end

Vagrant.configure('2') do |config|
  # We don't have NFS working inside Photon.
  config.nfs.functional = false

  %w('vmware_fusion', 'vmware_workstation', 'vmware_appcatalyst').each do |p|
    config.vm.provider p do |v|
      # Use paravirtualized virtual hardware on VMW hypervisors
      v.vmx['ethernet0.virtualDev'] = 'vmxnet3'
      v.vmx['scsi0.virtualDev'] = 'pvscsi'
    end
  end
end
