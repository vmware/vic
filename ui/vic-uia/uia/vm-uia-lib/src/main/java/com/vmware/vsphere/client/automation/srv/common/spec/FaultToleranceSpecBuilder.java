/* Copyright 2016 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.srv.common.spec;

import java.util.LinkedList;
import java.util.List;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import com.vmware.vsphere.client.automation.srv.common.spec.FaultToleranceSpec.FaultToleranceState;
import com.vmware.vsphere.client.automation.srv.common.spec.FaultToleranceSpec.FaultToleranceStatus;

/**
 * Class for creating Fault Tolerance Spec.
 */
public class FaultToleranceSpecBuilder {

   private FaultToleranceStatus ftStatus = FaultToleranceStatus.PROTECTED;
   private FaultToleranceState ftState = FaultToleranceState.STARTING;
   private List<String> tags = new LinkedList<>();

   private static final Logger _logger = LoggerFactory
         .getLogger(FaultToleranceSpecBuilder.class);

   public FaultToleranceSpecBuilder() {
      _logger.debug("Building fault tolerance spec:");
   }

   /**
    * Sets the fault tolerance status.
    *
    * @param ftStatus
    * @return FaultToleranceSpecBuilder
    */
   public FaultToleranceSpecBuilder setStatus(FaultToleranceStatus ftStatus) {
      this.ftStatus = ftStatus;
      _logger.debug(String.format("Fault tolerance status was set to %s",
            ftStatus.getValue()));
      return this;
   }

   /**
    * Sets the fault tolerance state.
    *
    * @param ftState
    * @return FaultToleranceSpecBuilder
    */
   public FaultToleranceSpecBuilder setState(FaultToleranceState ftState) {
      this.ftState = ftState;
      _logger.debug(String.format("Fault tolerance state was set to %s",
            ftState.getValue()));
      return this;
   }

   /**
    * Adds tags to the spec; used to differentiate multiple specs from one
    * another
    *
    * @param tags
    * @return FaultToleranceSpecBuilder
    */
   public FaultToleranceSpecBuilder setTags(String... tags) {
      for (String tag : tags) {
         this.tags.add(tag);
         _logger.debug(String.format("Assigned tag %s", tag));
      }
      return this;
   }

   /**
    * Builds a Fault Tolerance spec and returns it.
    *
    * @return FaultToleranceSpec
    */
   public FaultToleranceSpec build() {
      FaultToleranceSpec ftSpec = new FaultToleranceSpec();
      ftSpec.status.set(ftStatus);
      ftSpec.state.set(ftState);
      ftSpec.tag.set(tags.toArray(new String[] {}));
      return ftSpec;
   }
}