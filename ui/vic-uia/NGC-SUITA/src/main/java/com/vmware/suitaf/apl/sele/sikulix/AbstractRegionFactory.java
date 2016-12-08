/*
 * Copyright 2015 VMware, Inc. All rights reserved. -- VMware Confidential
 */
package com.vmware.suitaf.apl.sele.sikulix;

import org.sikuli.script.Region;

/**
 * Abstract factory used to create Sikuli {@link org.sikuli.script.Region}
 * objects.
 */
abstract class AbstractRegionFactory {
   /**
    * Implementations of this method return Region objects
    *
    * @return a new Sikuli region
    */
   public abstract Region getRegion();
}
