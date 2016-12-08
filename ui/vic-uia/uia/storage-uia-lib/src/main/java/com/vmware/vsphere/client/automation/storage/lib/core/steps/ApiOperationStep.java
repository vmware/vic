package com.vmware.vsphere.client.automation.storage.lib.core.steps;

import org.apache.commons.configuration.ConfigurationException;
import org.apache.commons.lang.NotImplementedException;

import com.vmware.client.automation.configuration.ConfigurationUtil;

/**
 * Base class for API operations
 */
public abstract class ApiOperationStep extends BasePrerequisiteStep {

   /**
    * Enumeration of the similarity of of objects
    */
   protected static enum Similarity {
      /**
       * The element is not match and no modification is going to change this
       */
      NO(0),
      /**
       * The element is similar, but additional actions are required in order to
       * be considered a match.
       */
      SIMILAR(1),
      /**
       * The element is considered a match without any modifications
       */
      MATCH(2);

      private final int similarityIndex;

      Similarity(int similarityIndex) {
         this.similarityIndex = similarityIndex;
      }
   }

   /**
    * Enumeration of all available Similarity thresholds
    */
   public enum SimilarityThreshold {
      /**
       * Threshold of disabled indicates that no similarity should be taken into
       * account
       */
      DISABLED(Integer.MAX_VALUE),

      /**
       * Indicates that both {@link Similarity} levels
       * {@link Similarity#SIMILAR} and {@link Similarity#MATCH} are taken into
       * account
       */
      MODIFICATION(Similarity.SIMILAR.similarityIndex),

      /**
       * Indicates that only {@link Similarity} of level
       * {@link Similarity#MATCH} is taken into account
       */
      MATCH(Similarity.MATCH.similarityIndex);

      private final int similarityIndex;

      SimilarityThreshold(int similarityIndex) {
         this.similarityIndex = similarityIndex;
      }
   }

   /**
    * Interface for execution of a clean operation
    */
   protected static interface CleanOperation {

      /**
       * The execution of the CleanOperation
       *
       * @throws Exception
       */
      public void execute() throws Exception;
   }

   private final SimilarityThreshold modificationThreshold;
   private Similarity similarity;
   private CleanOperation cleanOperation;

   public ApiOperationStep() {
      this(getDefaultModificationThreshold());
   }

   private static SimilarityThreshold getDefaultModificationThreshold() {
      try {
         return SimilarityThreshold.valueOf(ConfigurationUtil.getConfig()
               .getString("core.steps.apioperation.threshold",
                     SimilarityThreshold.DISABLED.name()));
      } catch (ConfigurationException e) {
         throw new RuntimeException(e);
      }
   }

   public ApiOperationStep(SimilarityThreshold modificationThreshold) {
      this.modificationThreshold = modificationThreshold;
   }

   @Override
   public final void execute() throws Exception {
      similarity = checkPresence();

      if (modificationThreshold.similarityIndex <= similarity.similarityIndex) {
         if (similarity == Similarity.SIMILAR) {
            cleanOperation = modify();
         }
      } else {
         cleanOperation = perform();
      }
   }

   /**
    * Implement to perform any modification for {@link Similarity#SIMILAR}
    * existence
    *
    * <p>
    * It is recommended to override this method if you override the
    * {@link ApiOperationStep#checkPresence()}
    * </p>
    *
    * @return
    * @throws Exception
    */
   protected CleanOperation modify() throws Exception {
      throw new NotImplementedException(String.format(
            "Implementation of modify is required for class %s.", getClass()
                  .getCanonicalName()));
   }

   /**
    * Implement a perform operation for the step. This method is used in case of
    * {@link Similarity#NO} or when {@link Similarity} is less than the
    * {@link SimilarityThreshold} for current istnace
    *
    * @return
    * @throws Exception
    */
   protected abstract CleanOperation perform() throws Exception;

   /**
    * Check for presence of the api operation outcome entity/configuration etc..
    *
    * <p>
    * It is recommended to override the {@link ApiOperationStep#modify()} if you
    * override this method
    * </p>
    *
    * <pre>
    * <b>Note:</b>
    * <p>This method is not guaranteed to be called when the Threshold is set to disabled </p>
    *
    * <pre>
    * @return
    */
   protected Similarity checkPresence() {
      return Similarity.NO;
   }

   @Override
   public final void clean() throws Exception {
      if (this.cleanOperation != null) {
         cleanOperation.execute();
      }
   }
}
