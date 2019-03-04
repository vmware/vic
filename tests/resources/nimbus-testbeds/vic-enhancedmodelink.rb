oneGB = 1 * 1000 * 1000 # in KB
nameParts = ['vic', 'em2-link', 'vcva']
$testbed = Proc.new do
  {
    'version' => 3,
    'name' => nameParts.join('-'),
  
      "esx" => (0..3).map do | idx |
        {
          "name" => "esx.#{idx}",
          "vc" => "vc.#{idx%2}",
          'cpus' => 8,
          'cpuReservation' => 4800,
          "style" => "fullInstall",
          "desiredPassword" => "ca$hc0w",
          "memory" => 32768, # 2x default
          'memoryReservation' => 16384,
          "disk" => [ 100 * oneGB, 100 * oneGB],
          "nics" => 2,
          "iScsi" => ["iscsi.#{idx%2+1}"],
          "clusterName" => "cls1",
          'localDatastoreNamePrefix' => "esx#{idx}-vmfs",
          'sharedDatastoreNamePrefix' => "sharedVmfs-",
        }
      end,
  
    "iscsi" => [
      {
        "name" => "iscsi.1",
        "luns" => [512],
        "iqnRandom" => "nimbus1"
      },
      {
        "name" => "iscsi.2",
        "luns" => [512],
        "iqnRandom" => "nimbus2"
      }
    ],
   
    'vcs' => [
       {
         'name' => "psc.0",
         'type' => 'vcva',
         'deploymentType' => 'infrastructure',
         'additionalScript' => [], # XXX: Create users
         'cpuReservation' => 2400,
         'memoryReservation' => 4096,
         'dbType' => 'embedded'
       },
       {
         'name' => "psc.1",
         'type' => 'vcva',
         'deploymentType' => 'infrastructure',
         'additionalScript' => [], # XXX: Create users
         'replicationPartner' => 'psc.0',
         'dependsOnHosts' => ['psc.0'],
         'cpuReservation' => 600,
         'memoryReservation' => 1200,
         'dbType' => 'embedded'
       },
       {
         'name' => "vc.0",
         'type' => 'vcva',
         'deploymentType' => 'management',
         'additionalScript' => [], # XXX: Create users
         'cpuReservation' => 600,
         'memoryReservation' => 1200,
         "dcName" => "dc1",
         "clusters" => [
           {"name" => "cls1", "vsan" => false, "enableDrs" => true, "enableHA" => true}
           ],
         "addHosts" => "allInSameCluster",
         'dirHost' => "psc.0",
         'dependsOnHosts' => ['psc.0'],
         'dbType' => 'embedded'
       },
       {
         'name' => "vc.1",
         'type' => 'vcva',
         'deploymentType' => 'management',
         'additionalScript' => [], # XXX: Create users
         'cpuReservation' => 600,
         'memoryReservation' => 1200,
         "dcName" => "dc1",
         "clusters" => [
           {"name" => "cls1", "vsan" => false, "enableDrs" => true, "enableHA" => true}
           ],
         "addHosts" => "allInSameCluster",
         'dirHost' => "psc.1",
         'dependsOnHosts' => ['psc.1'],
         'dbType' => 'embedded'
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
