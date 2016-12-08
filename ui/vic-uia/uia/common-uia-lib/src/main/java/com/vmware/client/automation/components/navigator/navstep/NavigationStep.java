/**
 * Copyright 2012 VMware, Inc. All rights reserved. -- VMware Confidential
 */
package com.vmware.client.automation.components.navigator.navstep;

import com.vmware.client.automation.components.navigator.spec.LocationSpec;

/**
 * The interface defines the API for one atomic navigation step.
 */
public interface NavigationStep {

   /**
    * The navigation ID associated with the navigation step.
    *
    * Specific tests should not know about a navigation step. They
    * should use this identifier to refer to it through the
    * <code>LocationSpec</code>.
    */
   public String getNId();

   /**
    * Use this method to provide an implementation of the actual UI
    * navigation step.
    *
    * The method shouldn't modify the state of the object.
    *
    * @param locationSpec
    *    A reference to the <code>LocationSpec</code>
    */
   public void doNavigate(LocationSpec locationSpec) throws Exception;
}

