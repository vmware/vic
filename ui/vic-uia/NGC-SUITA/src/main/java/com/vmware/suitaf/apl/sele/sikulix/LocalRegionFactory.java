/*
 * Copyright 2015 VMware, Inc. All rights reserved. -- VMware Confidential
 */
package com.vmware.suitaf.apl.sele.sikulix;

import org.sikuli.script.Region;
import org.sikuli.script.Sikulix;

/**
 * Implementation of the AbstractSikuliScreenFactory which creates and
 * returns local Sikuli Screen/Region objects
 */
class LocalRegionFactory extends AbstractRegionFactory {
   @Override
   public Region getRegion() {
      return Sikulix.init();
   }
}
