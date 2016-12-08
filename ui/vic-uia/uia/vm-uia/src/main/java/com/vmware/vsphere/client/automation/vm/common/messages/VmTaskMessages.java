package com.vmware.vsphere.client.automation.vm.common.messages;

import com.vmware.vsphere.client.test.i18n.gwt.Messages;

public interface VmTaskMessages extends Messages {
   @DefaultMessage("Power On virtual machine")
   String powerOn();

   @DefaultMessage("Power Off virtual machine")
   String powerOff();
}
