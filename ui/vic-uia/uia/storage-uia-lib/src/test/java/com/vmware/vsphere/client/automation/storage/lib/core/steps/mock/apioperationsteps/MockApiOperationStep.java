package com.vmware.vsphere.client.automation.storage.lib.core.steps.mock.apioperationsteps;

import com.vmware.vsphere.client.automation.storage.lib.core.steps.ApiOperationStep;

/**
 * Mock object for {@link ApiOperationStep}
 */
public class MockApiOperationStep extends ApiOperationStep {

   private boolean isCheckPresenceCalled = false;
   private boolean isPerformCalled = false;
   private boolean isModifyCalled = false;
   private boolean isPerformCleaned = false;
   private boolean isModifyCleaned = false;
   private final Similarity similarityToBeFound;

   /**
    * Initializes new instance of a {@link MockApiOperationStep}
    *
    * @param threshold
    *           the SimilarityThreshold for current instance. See
    *           {@link ApiOperationStep#ApiOperationStep(SimilarityThreshold)}
    * @param similarityToBeFound
    *           the result of {@link ApiOperationStep#checkPresence()}
    */
   public MockApiOperationStep(SimilarityThreshold threshold,
         Similarity similarityToBeFound) {
      super(threshold);
      this.similarityToBeFound = similarityToBeFound;
   }

   /**
    * Initializes new instance of a {@link MockApiOperationStep} with the
    * default {@link SimilarityThreshold}
    *
    * @param similarityToBeFound
    *           the result of {@link ApiOperationStep#checkPresence()}
    */
   public MockApiOperationStep(Similarity similarityToBeFound) {
      this.similarityToBeFound = similarityToBeFound;
   }

   @Override
   protected final CleanOperation perform() {
      this.isPerformCalled = true;
      return new CleanOperation() {
         @Override
         public void execute() {
            isPerformCleaned = true;
         }
      };
   }

   @Override
   protected CleanOperation modify() {
      this.isModifyCalled = true;
      return new CleanOperation() {
         @Override
         public void execute() {
            isModifyCleaned = true;
         }
      };
   }

   @Override
   protected Similarity checkPresence() {
      this.isCheckPresenceCalled = true;
      return this.similarityToBeFound;
   }

   public boolean isCheckPresenceCalled() {
      return isCheckPresenceCalled;
   }

   public boolean isPerformCalled() {
      return isPerformCalled;
   }

   public boolean isModifyCalled() {
      return isModifyCalled;
   }

   public boolean isPerformCleaned() {
      return isPerformCleaned;
   }

   public boolean isModifyCleaned() {
      return isModifyCleaned;
   }

}
