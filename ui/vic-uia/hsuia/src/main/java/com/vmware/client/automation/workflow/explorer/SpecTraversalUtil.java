/** Copyright 2014 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.client.automation.workflow.explorer;

import java.util.HashSet;
import java.util.LinkedHashSet;
import java.util.List;
import java.util.Set;

import org.apache.commons.collections4.CollectionUtils;

import com.vmware.client.automation.common.spec.BaseSpec;

/**
 * A set of utility methods for filtering and retrieving specific specs
 * from spec containers.
 */
public class SpecTraversalUtil {
   private SpecTraversalUtil() {
   }

   // Public utility methods =================================================================

   /**
    *
    * @param containerSpec
    * @param specClass
    * @return
    */
   public static <T extends BaseSpec> T getRequiredSpecFromContainerNode(
         final BaseSpec containerSpec, final Class <T> specClass)
               throws SpecNotFoundException, DuplicateSpecFoundException {

      final boolean traverseChildren = false;
      return getRequiredSpecFromContainerNode(containerSpec, specClass, traverseChildren);
   }

   //   public static <T extends BaseSpec> T getRequiredSpecFromContainerTree(
   //         final BaseSpec containerSpec, final Class <T> specClass)
   //               throws SpecNotFoundException, DuplicateSpecFoundException {
   //
   //      final boolean traverseChildren = true;
   //      return getRequiredSpecFromContainerNode(containerSpec, specClass, traverseChildren);
   //   }

   //   public static <T extends BaseSpec> Set<T> getAtLeastOneSpecFromContainerNode(
   //         final BaseSpec containerSpec, final Class<T> specClass) throws SpecNotFoundException  {
   //      validateInput(containerSpec, specClass);
   //
   //      final boolean traverseChildren = false;
   //      return getAtLeastOneFromContainer(containerSpec, specClass, traverseChildren);
   //   }
   //
   //   public static <T extends BaseSpec> Set<T> getAtLeastOneFromContainerTree(
   //         final BaseSpec containerSpec, final Class<T> specClass) throws SpecNotFoundException  {
   //
   //      final boolean traverseChildren = true;
   //      return getAtLeastOneFromContainer(containerSpec, specClass, traverseChildren);
   //   }

   public static <T extends BaseSpec> Set<T> getAllSpecsFromContainerNode(
         final BaseSpec containerSpec, final Class<T> specClass) {

      final boolean traverseChildren = false;
      return getAllSpecsFromContainer(containerSpec, specClass, traverseChildren);
   }

   public static <T extends BaseSpec> Set<T> getAllSpecsFromContainerTree(
         final BaseSpec containerSpec, final Class<T> specClass) {

      final boolean traverseChildren = true;
      return getAllSpecsFromContainer(containerSpec, specClass, traverseChildren);
   }

   // Helper methods =================================================================================

   private static <T extends BaseSpec> T getRequiredSpecFromContainerNode(
         final BaseSpec containerSpec, final Class <T> specClass,
         final boolean traverseChildren) throws SpecNotFoundException, DuplicateSpecFoundException {
      validateInput(containerSpec, specClass);

      final Set<T> retrievedSpecs = new LinkedHashSet<T>(); // LinkedHashSet - keep the order.
      traverseSpecLink(retrievedSpecs, containerSpec, specClass, traverseChildren);

      if (retrievedSpecs.isEmpty()) {
         throw new SpecNotFoundException(
               String.format(
                     "Cannot find the required spec class: %s.",
                     specClass.getCanonicalName()));
      }

      if (retrievedSpecs.size() > 1) {
         throw new DuplicateSpecFoundException(
               String.format(
                     "More than single spec of type %s was found, but single one is expected.",
                     specClass.getCanonicalName()));
      }

      return retrievedSpecs.iterator().next();
   }

   //   private static <T extends BaseSpec> Set<T> getAtLeastOneFromContainer(
   //         final BaseSpec containerSpec, final Class<T> specClass,
   //         final boolean traverseChildren) throws SpecNotFoundException  {
   //
   //      validateInput(containerSpec, specClass);
   //
   //      final Set<T> retrievedSpecs = new LinkedHashSet<T>(); // LinkedHashSet - keep the order.
   //      traverseSpecLink(retrievedSpecs, containerSpec, specClass, traverseChildren);
   //
   //      if (retrievedSpecs.isEmpty()) {
   //         throw new SpecNotFoundException(
   //               String.format(
   //                     "Cannot find any instance from the required spec class: %s. At lest one is required",
   //                     specClass.getCanonicalName()));
   //      }
   //
   //      return retrievedSpecs;
   //   }


   private static <T extends BaseSpec> Set<T> getAllSpecsFromContainer(
         final BaseSpec containerSpec, final Class<T> specClass, final boolean traverseChildren) {

      validateInput(containerSpec, specClass);

      final Set<T> retrievedSpecs = new LinkedHashSet<T>(); // LinkedHashSet - keep the order.
      traverseSpecLink(retrievedSpecs, containerSpec, specClass, traverseChildren);

      return retrievedSpecs;
   }

   private static <T extends BaseSpec> void  traverseSpecLink(
         final Set<T> retrievedSpecs, final BaseSpec containerSpec,
         final Class<T> specClass, final boolean traverseChildren) {

      final Set<BaseSpec> processedSpecs = new HashSet<BaseSpec>();

      traverseSpecLink(retrievedSpecs, containerSpec, processedSpecs, specClass, traverseChildren);
   }

   @SuppressWarnings("unchecked")
   private static <T extends BaseSpec> void traverseSpecLink(
         final Set<T> retrievedSpecs, final BaseSpec containerSpec,
         final Set<BaseSpec> processedSpecs, final Class<T> specClass,
         final boolean traverseChildren) {

      final List<BaseSpec> linkedSpecs = containerSpec.links.getAll(BaseSpec.class);
      if (CollectionUtils.isEmpty(linkedSpecs)) {
         return;
      }

      for (BaseSpec spec : linkedSpecs) {
         if (!processedSpecs.contains(spec)) {
            processedSpecs.add(spec);
            // if (spec.getClass().isAssignableFrom(specType)) {
            if (specClass.isAssignableFrom(spec.getClass())) {
               retrievedSpecs.add((T)spec);
            }
            if (traverseChildren) {
               traverseSpecLink(
                     retrievedSpecs, spec /*new container*/, processedSpecs, specClass, traverseChildren);
            }
         }
      }
   }

   /**
    * Performs basic validation of the input parameters.
    *
    * Throws IllegalArgumentException of either of the parameters is not supplied.
    */
   private static void validateInput(
         final BaseSpec containerSpec, final Class<? extends BaseSpec> specClass) {
      if (containerSpec == null) {
         throw new IllegalArgumentException(
               "Required input paramater containerSpec not set.");
      }

      if (specClass == null) {
         throw new IllegalArgumentException(
               "Required input parameter specClass not set.");
      }
   }
}
