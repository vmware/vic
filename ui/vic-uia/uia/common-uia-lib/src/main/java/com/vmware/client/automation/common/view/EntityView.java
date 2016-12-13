/* Copyright 2015 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.client.automation.common.view;

import com.vmware.client.automation.util.UiDelay;

/**
 * Implements the main view for an entity.
 */
public class EntityView extends BaseView {

   private static final String ID_ENTITY_NAME_LABEL = "title";

   /**
    * Gets the name of the focused entity.
    *
    * @return name of the entity
    */
   public String getEntityName() {
      UI.condition.isFound(ID_ENTITY_NAME_LABEL).await(
         UiDelay.PAGE_LOAD_TIMEOUT.getDuration());

      return UI.component.value.get(ID_ENTITY_NAME_LABEL);
   }
}
