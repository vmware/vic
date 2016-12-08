/* Copyright 2013 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.srv.common.step;

import java.util.ArrayList;
import java.util.List;

import com.vmware.client.automation.workflow.BaseWorkflowStep;
import com.vmware.vsphere.client.automation.srv.common.spec.FolderSpec;
import com.vmware.vsphere.client.automation.srv.common.spec.FolderType;
import com.vmware.vsphere.client.automation.srv.common.srvapi.FolderBasicSrvApi;

public class CreateDatacenterFolderStep extends BaseWorkflowStep {

   private List<FolderSpec> _foldersToCreate;
   private List<FolderSpec> _foldersToDelete;

   @Override
   /**
    * @inheritDoc
    */
   public void prepare() throws Exception {

      List<FolderSpec> allFolderSpecs = getSpec().links
            .getAll(FolderSpec.class);

      if (allFolderSpecs == null || allFolderSpecs.size() == 0) {
         throw new IllegalArgumentException(
               "The spec has no links to 'FolderSpec' instances");
      }

      _foldersToCreate = getDcFolderSpecs(allFolderSpecs);

      if (_foldersToCreate == null || _foldersToCreate.size() == 0) {
         throw new IllegalArgumentException(
               "The spec has no links to DC folder 'FolderSpec' instances");
      }

      _foldersToDelete = new ArrayList<FolderSpec>();
   }

   @Override
   public void execute() throws Exception {

      for (FolderSpec folderSpec : _foldersToCreate) {
         if (!FolderBasicSrvApi.getInstance().createFolder(folderSpec)) {
            throw new Exception(String.format("Unable to create folder '%s'",
                  folderSpec.name.get()));
         }
         _foldersToDelete.add(folderSpec);
      }
   }

   @Override
   /**
    * @inheritDoc
    */
   public void clean() throws Exception {
      for (FolderSpec folderSpec : _foldersToCreate) {
         FolderBasicSrvApi.getInstance().deleteFolderSafely(folderSpec);
      }
   }

   // Extracts the folder specs that have DC type set
   private List<FolderSpec> getDcFolderSpecs(List<FolderSpec> folderSpecs) {
      List<FolderSpec> dcFolderSpecs = new ArrayList<FolderSpec>();
      for (FolderSpec folderSpec : folderSpecs) {
         if (folderSpec.type.get() == FolderType.DATACENTER) {
            dcFolderSpecs.add(folderSpec);
         }
      }

      return dcFolderSpecs.size() == 0 ? null : dcFolderSpecs;
   }

}
