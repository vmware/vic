/** Copyright 2014 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.client.automation.workflow.common;

/**
 * Defines workflow step context. The step context keeps the state of the step
 * and reference to the step itself. this base class provides common
 * functionality for each step - access methods for title, tag filer and etc.
 */
public class WorkflowStepContext {
   private WorkflowStep _step;

   private String _title = "";
   private String[] _specTagsFilter;

   /**
    * Getter for the context step.
    * @return
    */
   public WorkflowStep getStep() {
      return this._step;
   }

   /**
    * Getter for the array of assigned step tags.
    * @return
    */
   public String[] getSpecTagsFilter() {
      return _specTagsFilter;
   }

   /**
    * Setter for the assigned step tags.
    * @param specTagsFilter
    */
   public void setSpecTagsFilter(String[] specTagsFilter) {
      this._specTagsFilter = specTagsFilter;
   }

   /**
    * Getter for the step title.
    * @return
    */
   public String getTitle() {
      return _title;
   }

   /**
    * Setter for the step title.
    * @param title
    */
   public void setTitle(String title) {
      this._title = title;
   }

}
