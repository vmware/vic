package com.vmware.vsphere.client.automation.storage.lib.core.steps;

import junit.framework.Assert;

import org.junit.Test;

import com.vmware.vsphere.client.automation.storage.lib.core.steps.ApiOperationStep.Similarity;
import com.vmware.vsphere.client.automation.storage.lib.core.steps.ApiOperationStep.SimilarityThreshold;
import com.vmware.vsphere.client.automation.storage.lib.core.steps.mock.apioperationsteps.MockApiOperationStep;

public class ApiOperationStepTests {

   /**
    * Process a {@link MockApiOperationStep} instance with a given similarity
    * threshold and a given similarity to be found
    *
    * @param threshold
    * @param similarityToBeFound
    * @return
    */
   private static MockApiOperationStep processStep(
         SimilarityThreshold threshold, Similarity similarityToBeFound) {
      MockApiOperationStep operationStep = new MockApiOperationStep(threshold,
            similarityToBeFound);
      return processStep(operationStep);
   }

   /**
    * Process a {@link MockApiOperationStep}
    *
    * @param operationStep
    * @return
    */
   private static MockApiOperationStep processStep(
         MockApiOperationStep operationStep) {
      try {
         operationStep.prepare();
         operationStep.execute();
         operationStep.clean();
      } catch (Exception e) {
         throw new RuntimeException(e);
      }

      return operationStep;
   }

   @Test
   public void defaultThresholdSimilarityNo() {
      MockApiOperationStep noMatchStep = processStep(new MockApiOperationStep(
            Similarity.NO));

      verifyOperationStepPerformed(noMatchStep);
   }

   @Test
   public void defaultThresholdSimilaritySimilar() {
      MockApiOperationStep similarStep = processStep(new MockApiOperationStep(
            Similarity.SIMILAR));

      verifyOperationStepPerformed(similarStep);
   }

   @Test
   public void defaultThresholdSimilarityMatch() {
      MockApiOperationStep matchStep = processStep(new MockApiOperationStep(
            Similarity.MATCH));

      verifyOperationStepPerformed(matchStep);
   }

   @Test
   public void thresholdDisabledSimilarityNo() {
      MockApiOperationStep noMatchStep = processStep(
            SimilarityThreshold.DISABLED, Similarity.NO);

      verifyOperationStepPerformed(noMatchStep);
   }

   @Test
   public void thresholdDisabledSimilaritySimilar() {
      MockApiOperationStep similarStep = processStep(
            SimilarityThreshold.DISABLED, Similarity.SIMILAR);

      verifyOperationStepPerformed(similarStep);
   }

   @Test
   public void thresholdDisabledSimilarityMatch() {
      MockApiOperationStep matchStep = processStep(
            SimilarityThreshold.DISABLED, Similarity.MATCH);

      verifyOperationStepPerformed(matchStep);
   }

   @Test
   public void thresholdModificationSimilarityNo() {
      MockApiOperationStep noMatchStep = processStep(
            SimilarityThreshold.MODIFICATION, Similarity.NO);

      verifyOperationStepPerformed(noMatchStep);
   }

   @Test
   public void thresholdModificationSimilaritySimilar() {
      MockApiOperationStep similarStep = processStep(
            SimilarityThreshold.MODIFICATION, Similarity.SIMILAR);

      verifyOperationStepModified(similarStep);
   }

   @Test
   public void thresholdModificationSimilarityMatch() {
      MockApiOperationStep matchStep = processStep(
            SimilarityThreshold.MODIFICATION, Similarity.MATCH);

      verifyOperationStepNoActionsTaken(matchStep);
   }

   @Test
   public void thresholMatchSimilarityNo() {
      MockApiOperationStep matchStep = processStep(SimilarityThreshold.MATCH,
            Similarity.NO);

      verifyOperationStepPerformed(matchStep);
   }

   @Test
   public void thresholdMatchSimilaritySimilar() {
      MockApiOperationStep matchStep = processStep(SimilarityThreshold.MATCH,
            Similarity.SIMILAR);

      verifyOperationStepPerformed(matchStep);
   }

   @Test
   public void thresholdMatchSimilarityMatch() {
      MockApiOperationStep matchStep = processStep(SimilarityThreshold.MATCH,
            Similarity.MATCH);

      verifyOperationStepNoActionsTaken(matchStep);
   }

   /**
    * Verify the operation step in state that Perform action took place
    *
    * @param matchStep
    */
   private static void verifyOperationStepPerformed(
         MockApiOperationStep operationStep) {
      Assert.assertEquals("Perform is trigered", true,
            operationStep.isPerformCalled());
      Assert.assertEquals("Modifycation *NOT* trigered", false,
            operationStep.isModifyCalled());

      Assert.assertEquals("Perform is cleaned", true,
            operationStep.isPerformCleaned());
      Assert.assertEquals("Modifycation is *NOT* cleaned", false,
            operationStep.isModifyCleaned());
   }

   /**
    * Verify the operation step in state that Modify action took place
    *
    * @param matchStep
    */
   private static void verifyOperationStepModified(
         MockApiOperationStep operationStep) {
      Assert.assertEquals("Perform is *NOT* trigered", false,
            operationStep.isPerformCalled());
      Assert.assertEquals("Modifycation trigered", true,
            operationStep.isModifyCalled());

      Assert.assertEquals("Perform is *NOT* cleaned", false,
            operationStep.isPerformCleaned());
      Assert.assertEquals("Modifycation is cleaned", true,
            operationStep.isModifyCleaned());
   }

   /**
    * Verify the operation step in state that no actions are taken over the step
    *
    * @param operationStep
    */
   private static void verifyOperationStepNoActionsTaken(
         MockApiOperationStep operationStep) {
      Assert.assertEquals("Perform is *NOT* trigered", false,
            operationStep.isPerformCalled());
      Assert.assertEquals("Modifycation *NOT* trigered", false,
            operationStep.isModifyCalled());

      Assert.assertEquals("Perform is *NOT* cleaned", false,
            operationStep.isPerformCleaned());
      Assert.assertEquals("Modifycation is *NOT* cleaned", false,
            operationStep.isModifyCleaned());
   }

}
