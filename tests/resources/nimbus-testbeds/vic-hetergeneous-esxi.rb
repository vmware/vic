oneGB = 1 * 1000 * 1000 # in KB
 
$testbed = Proc.new do
  {
    "name" => "vic-hetegeneous",
    "version" => 3,
    "esx" => [
      {
        "name" => "esx.0",
        "vc" => "vc.0",
        'cpus' => 8,
        'cpuReservation' => 2400,
        "style" => "fullInstall",
        "desiredPassword" => "e2eFunctionalTest",
        "memory" => 16384, # 2x default
        'memoryReservation' => 8192,
        "disk" => [ 30 * oneGB],
        "nics" => 2,
        "iScsi" => ["iscsi.0"],
        "clusterName" => "cls1",
      },
      {
        "name" => "esx.1",
        "vc" => "vc.0",
        'cpus' => 8,
        'cpuReservation' => 2400,
        "style" => "fullInstall",
        "desiredPassword" => "e2eFunctionalTest",
        "memory" => 16384, # 2x default
        'memoryReservation' => 8192,
        "disk" => [ 30 * oneGB],
        "nics" => 2,
        "iScsi" => ["iscsi.0"],
        "clusterName" => "cls1",
        "customBuild" => "8935087"  #65U2
      },
      {
        "name" => "esx.2",
        "vc" => "vc.0",
        'cpus' => 8,
        'cpuReservation' => 2400,
        "style" => "fullInstall",
        "desiredPassword" => "e2eFunctionalTest",
        "memory" => 16384, # 2x default
        'memoryReservation' => 8192,
        "disk" => [ 30 * oneGB],
        "nics" => 2,
        "iScsi" => ["iscsi.0"],
        "clusterName" => "cls1",
        "customBuild" => "5050593"  #60
      }
    ],

    "iscsi" => [
      {
        "name" => "iscsi.0",
        "luns" => [200],
        "iqnRandom" => "nimbus1"
      }
    ],

    "vcs" => [
      {
        "name" => "vc.0",
        "type" => "vcva",
        'cpuReservation' => 2400,
        'memoryReservation' => 4096,
        "dcName" => "dc1",
        "clusters" => [{"name" => "cls1", "vsan" => false, "enableDrs" => true, "enableHA" => true}],
        "addHosts" => "allInSameCluster",
      }
    ],

    "postBoot" => Proc.new do |runId, testbedSpec, vmList, catApi, logDir|
      esxList = vmList['esx']
        esxList.each do |host|
          host.ssh do |ssh|
            ssh.exec!("esxcli network firewall set -e false")
          end
      end
    end
  }
end

