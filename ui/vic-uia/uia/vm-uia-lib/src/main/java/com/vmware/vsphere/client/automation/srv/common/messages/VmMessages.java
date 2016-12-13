package com.vmware.vsphere.client.automation.srv.common.messages;

import com.vmware.vsphere.client.test.i18n.gwt.Messages;

public interface VmMessages extends Messages {
   @DefaultMessage("Name")
   String getNameColumnHeader();

   @DefaultMessage("Contents")
   String getContentsColumnHeader();
}
