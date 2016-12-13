/** Copyright 2014 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.client.automation.workflow.provider;

import java.util.AbstractMap.SimpleEntry;
import java.util.List;

import com.vmware.client.automation.common.spec.EntitySpec;
import com.vmware.client.automation.workflow.common.WorkflowSpec;
import com.vmware.client.automation.workflow.explorer.TestbedSpecConsumer;
import com.vmware.client.automation.workflow.explorer.TestbedSpecProvider;
import com.vmware.hsua.common.datamodel.BasePBox.DataProperty;

/**
 * Container for the published by the providers specs.
 */
public class PublisherSpec extends WorkflowSpec implements TestbedSpecConsumer,
      TestbedSpecProvider {
   public DataProperty<SimpleEntry<String, EntitySpec>> entitySpecMap;

   @Override
   public void publishEntitySpec(String entityId, EntitySpec entitySpec) {
      List<SimpleEntry<String, EntitySpec>> list = entitySpecMap.getAll();

      EntitySpec existingSpec = getEntitySpecById(list, entityId);

      if (existingSpec != null) {
         throw new IllegalArgumentException("Duplicate");
      }

      list.add(new SimpleEntry<String, EntitySpec>(entityId, entitySpec));
      entitySpecMap.set(list);
   }

   @Override
   public <T extends EntitySpec> T getPublishedEntitySpec(String entityId) {
      List<SimpleEntry<String, EntitySpec>> list = entitySpecMap.getAll();
      EntitySpec existingSpec = getEntitySpecById(list, entityId);

      if (existingSpec == null) {
         throw new IllegalArgumentException("No spec founds for entityId: " + entityId);
      }
      @SuppressWarnings("unchecked")
      T t = (T) existingSpec;
      return t;
   }

   private EntitySpec getEntitySpecById(List<SimpleEntry<String, EntitySpec>> list,
         String entityId) {
      for (SimpleEntry<String, EntitySpec> entry : list) {
         if (entry.getKey().equals(entityId)) {
            return entry.getValue();
         }
      }

      return null;
   }
}
