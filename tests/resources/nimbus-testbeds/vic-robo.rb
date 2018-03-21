#usage: nimbus-testbeddeploy --testbedSpecRubyFile ~/vcan-vc-5esx-allflash.rb \
# --testbedName vcan-vc-5esx-allflash-iscsi-fullInstall-vcva --vcvaBuild ob-4602587 \
# --esxBuild ob-3620759 --runName NAME --esxPxeDir ob-3620759
#
# This file is needed because there are 3 functions in it which are used in this testbed.
require '/mts/git/nimbus/lib/testframeworks/testng/testng.rb'
oneGB = 1 * 1000 * 1000
$testbed = Proc.new do |sharedStorageStyle, esxStyle, vcStyle|
  esxStyle ||= 'fullInstall'
  vcStyle ||= 'vcva'
  esxStyle = esxStyle.to_s
  vcStyle = vcStyle.to_s
  sharedStorageStyle = sharedStorageStyle.to_s

  testbed = {
    'version' => 3,
    'name' => "vcan-vc-3esx-allflash-#{sharedStorageStyle}-#{esxStyle}-#{vcStyle}",
    'esx' => (0..2).map do |i|
      {
        'name' => "esx.#{i}",
        'style' => esxStyle,
        'numMem' => 64 * 1024,
        'numCPUs' => 6,
        'disks' => [ 250 * oneGB ],
        'disableNfsMounts' => true,
        'ssds' => [ 200 * oneGB, 250 * oneGB ],
        'nics' => 2,
        'staf' => false,
        'desiredPassword' => 'ca$hc0w',
        'vmotionNics' => ['vmk0'],
        'clusterName' => '001',
        'post_postboot_commands' => [
            'esxcli system settings kernel set -s autoPartition -v "TRUE"', \
            'esxcli system settings kernel set -s skipPartitioningSsds -v "TRUE"', \
            'services.sh restart', \
            'esxcli vsan storage tag add -d mpx.vmhba1:C0:T1:L0 -t capacityFlash', \
            'esxcli storage core adapter rescan --all', \
            'esxcli network firewall set --enabled false', \
            'echo \'vhv.enable = "TRUE"\' >> /etc/vmware/config'
            ],
      }
    end,
    'vcs' => [{
      'name' => 'vc',
      'type' => vcStyle,
      'numMem' => 32 * 1024,
      'addHosts' => 'allInSameCluster',
      'clusters' => [{'name' => '001', 'dc' => 'vcqaDC', 'vsan' => true, 'enableDrs' => false, 'enableHA' => true}],
      'linuxCertFile' => TestngLauncher.vcvaCertFile,
#     'addHosts' => 'multiple-clusters',
#      'esxsPerCluster' => '3',
#      'numOfClusters' => '2',
#      'clusterName' => ['001', '002'],
             }],
      'vsan' => true,
    'postBoot' => Proc.new do |runId, testbedSpec, vmList|
      vc = vmList['vc'].first
      puts "======#{vc.ip}======"
    end,
  }

  testbed
end

[:pxeBoot, :fullInstall].each do |esxStyle|
  [:vpxInstall, :vcva].each do |vcStyle|
    [:iscsi, :fc].each do |sharedStorageStyle|
      testbedSpec = $testbed.call(sharedStorageStyle, esxStyle, vcStyle)
      Nimbus::TestbedRegistry.registerTestbed testbedSpec
    end
  end
end
