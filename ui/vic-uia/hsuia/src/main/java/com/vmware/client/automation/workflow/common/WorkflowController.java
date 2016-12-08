/** Copyright 2014 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.client.automation.workflow.common;

import java.lang.reflect.InvocationTargetException;
import java.util.ArrayList;
import java.util.List;

import com.vmware.client.automation.common.spec.BaseSpec;
import com.vmware.client.automation.workflow.explorer.WorkflowRegistry;
import com.vmware.client.automation.workflow.provider.ProviderControllerException;

/**
 * Base class for workflow controllers. Contains common methods for work with
 * specs.
 */
public abstract class WorkflowController {
   private boolean _isScenarioLaunched = false;

   public abstract WorkflowRegistry getRegistry();

   /**
    * Make a deep clone of a spec including the linked specs. The routine will
    * also apply tags filter, if specified, cloning only specs that have at
    * least one matching tag with the filter. Tag filtering is not performed on
    * the root spec, which is the workflow spec.
    * 
    * @param sourceSpec
    *           A reference to the spec
    * @param tagsFilter
    *           An array of tags to use as filter
    * @return Cloned spec.
    */
   protected BaseSpec getDeepClonedSpec(BaseSpec sourceSpec, String[] tagsFilter)
         throws IllegalArgumentException, SecurityException, InstantiationException,
         IllegalAccessException, InvocationTargetException, NoSuchMethodException {

      return getDeepClonedSpecInternal(sourceSpec, tagsFilter, new ArrayList<BaseSpec>());
   }


   /**
    * Makes sure that only a single request for a single scenario has been made
    * to the controller.
    * 
    * The workflow controller is generally designed to work with the context in
    * a way that it can support running only one scenario at a time.
    * 
    * @throws ProviderControllerException
    *            This exception is thrown in case previous request to a scenario
    *            has been made.
    */
   protected void guardSingleRequest() throws ProviderControllerException {
      if (_isScenarioLaunched) {
         throw new ProviderControllerException(
               "The controller has been already used to run provider scenarion."
                     + " Only single run is supported. Release the controller.",
                     null /* This is called in that way by exception */);
      } else {
         _isScenarioLaunched = true;
      }
   }


   /**
    * A helper method that implements the spec cloning as described in the
    * <code>getDeepClonedSpec</code> routine.
    * 
    * @param sourceSpec
    *           A reference to the spec
    * @param tagsFilter
    *           An array of tags to use as filter
    * @param processedSourceSpecs
    *           A list of the already cloned source specs. The
    *           <code>PropertyBox</code> links automatically provide references
    *           to all the parent specs, which should be excluded from this
    *           cloning.
    * 
    * @return Cloned spec.
    * @throws IllegalAccessException
    * @throws InstantiationException
    */
   private BaseSpec getDeepClonedSpecInternal(BaseSpec sourceSpec, String[] tagsFilter,
         List<BaseSpec> processedSourceSpecs) throws InstantiationException,
         IllegalAccessException {

      // Empty spec. there's nothing to process.
      if (sourceSpec == null) {
         return null;
      }

      // Skip already processed specs. This is necessary because spec links
      // provide automatically references to their container specs.
      if (processedSourceSpecs.contains(sourceSpec)) {
         return null;
      }

      boolean cloneSpec = true;

      // Try to match tags if it's not the workflow container assuming the root
      // is always the workflow container.
      if (processedSourceSpecs.size() > 0 && requireTagMatchEval(tagsFilter)
            && !hasTagMatch(sourceSpec, tagsFilter)) {
         cloneSpec = false;
      }

      processedSourceSpecs.add(sourceSpec);

      BaseSpec newSpec = null;

      if (cloneSpec) {
         // Initialize the new spec from the correct type
         newSpec = sourceSpec.getClass().newInstance();

         // Copy the data properties
         newSpec.copy(sourceSpec);

         // Copy the linked specs
         List<BaseSpec> linkedSourceSpecs = sourceSpec.links.getAll(BaseSpec.class);
         for (BaseSpec linkedSourceSpec : linkedSourceSpecs) {

            BaseSpec newLinkedSpec =
                  getDeepClonedSpecInternal(
                        linkedSourceSpec,
                        tagsFilter,
                        processedSourceSpecs);

            if (newLinkedSpec != null) {
               newSpec.links.add(newLinkedSpec);
            }
         }
      }

      return newSpec;
   }

   /**
    * Check if at least one of the tags in the filter matches a tag in the spec.
    * 
    * @param sourceSpec
    *           A spec reference
    * @param tagsFilter
    *           An array of tags to use as filter.
    * 
    * @return true - at least one filter tag has a match in the spec; false -
    *         otherwise;
    */
   private boolean hasTagMatch(BaseSpec spec, String[] tagsFilter) {
      // Check if the spec has the required tags and it's not the workflow spec.
      if (spec.tag.isAssigned()) {
         boolean tagFound = false;

         List<String> specTags = spec.tag.getAll();
         for (String tagFilter : tagsFilter) {
            if (specTags.contains(tagFilter)) {
               tagFound = true;
               break;
            }
         }

         return tagFound;
      }

      return false;
   }

   /**
    * Check if evaluation for matching tags between the spec and the tags filter
    * should be performed.
    * 
    * @param tagsFilter
    *           An array of tags to use as filter
    * 
    * @return true - the tags filter is not empty; false - otherwise;
    */
   private boolean requireTagMatchEval(String[] tagsFilter) {
      return tagsFilter != null && tagsFilter.length > 0;
   }
}
