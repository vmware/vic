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

       vcs = vmList['vc']
       vcs.each do |vc|
         ip = vc.ip
         name = vc.name
         type = vc.deploymentType
         Log.info  "#{name}, #{ip}, #{type}"
         if type == 'infrastructure' 
           next
         end
         Log.info "Configure on #{name}"
          
       #end
       #vcs = vmList['vc']
       #[2..3].map do |id|
         
         #vc = vcs[id]
         #ip = vc.ip
         Log.info  "#{ip}"
         vim = VIM.connect vc.rbvmomiConnectSpec
         datacenters = vim.serviceInstance.content.rootFolder.childEntity.grep(RbVmomi::VIM::Datacenter)
         raise "Couldn't find a Datacenter precreated"  if datacenters.length == 0
         datacenter = datacenters.first
         Log.info "Found a datacenter successfully in the system, name: #{datacenter.name}"
         clusters = datacenter.hostFolder.children
         raise "Couldn't find a cluster precreated"  if clusters.length == 0
         cluster = clusters.first
         Log.info "Found a cluster successfully in the system, name: #{cluster.name}"
   
         dvs = datacenter.networkFolder.CreateDVS_Task(
           :spec => {
             :configSpec => {
               :name => "test-ds"
             },
   	    }
         ).wait_for_completion
         Log.info "Vds DSwitch created"
   
         dvpg1 = dvs.AddDVPortgroup_Task(
           :spec => [
             {
               :name => "management",
               :type => :earlyBinding,
               :numPorts => 12,
             }
           ]
         ).wait_for_completion
         Log.info "management DPG created"

         dvpg2 = dvs.AddDVPortgroup_Task(
           :spec => [
             {
               :name => "vm-network",
               :type => :earlyBinding,
               :numPorts => 12,
             }
           ]
         ).wait_for_completion
         Log.info "vm-network DPG created"
   
         dvpg3 = dvs.AddDVPortgroup_Task(
           :spec => [
             {
               :name => "bridge",
               :type => :earlyBinding,
               :numPorts => 12,
             }
           ]
         ).wait_for_completion
         Log.info "bridge DPG created"
   
         Log.info "Add hosts to the DVS"
         onecluster_pnic_spec = [ VIM::DistributedVirtualSwitchHostMemberPnicSpec({:pnicDevice => 'vmnic1'}) ]
         dvs_config = VIM::DVSConfigSpec({
           :configVersion => dvs.config.configVersion,
           :host => cluster.host.map do |host|
           {
             :operation => :add,
             :host => host,
             :backing => VIM::DistributedVirtualSwitchHostMemberPnicBacking({
               :pnicSpec => onecluster_pnic_spec
             })
           }
           end
         })
         dvs.ReconfigureDvs_Task(:spec => dvs_config).wait_for_completion
         Log.info "Hosts added to DVS successfully"
       end
     end
  }
 end
