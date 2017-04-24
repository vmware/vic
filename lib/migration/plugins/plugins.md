# Developer guideline to write a data migration plugin

## Who and When to develop a data migration plugin

As [upgrade design](../../../doc/design/upgrade.md) mentioned, every developer who changed configuration from VCH appliance configuration, container vm configuration, or key/value store, whatever it is, they should add data migration plugin here.

## How to develop a data migration plugin

### VM guestinfo changes including VCH appliance and container VM

- Be sure to check if there is file conflict with other's commit. If there is conflict, that means others had plugin check in as well, make sure you did not use same version number with others
- Create one package in this directory for your plugin
- Add plugin files to the new package
- Your plugin might rely on previous version and current version, if you need the whole VCH appliance or container configuration definition, copy them to package lib/migration/config, with correct new package created. For example, if you are working on migration from version 4 to version 5, create a package named v4 if there is no one existing, and copy all old configuration files to there, and create a new package v5,  copy all new configuration files to v5 package. Remember to update package name.
- If only a few attributes are touched, you could define that piece of attribute definition in your plugin package, without whole configuration files copied in the above step. Which can save binary size, and also configuration encodeing/decoding time.
- Develop plugin based on your specific change, read the input data, decoding to acceptable format, and change value inside of them. stop_singal_rename_sample.go is one VCH appliance configuration change sample, you can follow that to write your own plugin.
- register your plugin to data migration framework, and put the function in package init() method.
  * Register appliance data migration plugin to manager.ApplianceConfigure category
  * Register container data migration plugin to manager.ContainerConfigure category
  * The two kinds of plugs target for different data, be sure the type is selected correctly in registration.
- Add import of your package in init.go in this package, to make sure the plugin is registered dynamically.

### key/value store changes

- Be sure to check if there is file conflict with other's commit. If there is conflict, that means others had plugin check in as well, make sure you did not use same version number with others
- Create one package in this directory for your plugin
- Add your plugin files to the new package
- Develop plugin based on your specific change.
  * Create new key/value store file in datastore, instead of change in existing file. The new file for new version, e.g. v4, could use file name as XXX.v4, to differentiate with old version's file.
  * Copy all existing key/values to the new file, and update to the new file directly.
  * Write to new datastore file only, because old file is still correct if migration failed eventually.
- register your plugin to data migration framework, and put the function in package init() method.
  * Register key/value store data migration plugin to manager.ApplianceConfigure category
- Add import of your package in init.go in this package, to make sure the plugin is registered dynamically.

Note:
- Plugin version should be greater than 0
- If you changed both key/value store and VCH appliance configuration, please add two separate plugins for them. Eventually, data migration framework will execute both of them, but separation will make the code easy to read.
- While copy configuration files, remove unnecessary methods from that file, to reduce binary file size

## Add integration test

Be sure to add test scenario into upgrade test group, to cover your changes, to make sure after upgrade, your function works well
