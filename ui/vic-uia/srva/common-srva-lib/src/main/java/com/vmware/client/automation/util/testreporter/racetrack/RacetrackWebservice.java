package com.vmware.client.automation.util.testreporter.racetrack;

import java.io.File;
import java.io.FileNotFoundException;
import java.io.IOException;
import java.net.MalformedURLException;
import java.net.URISyntaxException;
import java.net.URL;
import java.util.ArrayList;
import java.util.Enumeration;
import java.util.Hashtable;
import java.util.List;

import org.apache.commons.httpclient.DefaultHttpMethodRetryHandler;
import org.apache.commons.httpclient.HttpClient;
import org.apache.commons.httpclient.HttpException;
import org.apache.commons.httpclient.HttpStatus;
import org.apache.commons.httpclient.NameValuePair;
import org.apache.commons.httpclient.methods.PostMethod;
import org.apache.commons.httpclient.methods.multipart.FilePart;
import org.apache.commons.httpclient.methods.multipart.MultipartRequestEntity;
import org.apache.commons.httpclient.methods.multipart.Part;
import org.apache.commons.httpclient.methods.multipart.StringPart;
import org.apache.commons.httpclient.params.HttpMethodParams;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.testng.TestException;

/**
 * This class encapsulates the functionality to interact with the
 * Racetrack web service.
 *
 * @author gautam
 * @author michaely
 *
 * Use this singleton in either of the following two approaches:
 *
 * 1. Recommended way of using this Singleton: always call getInstance()
 * with no arguments. The default racetrackUrl is "racetrack.eng.vmware.com".
 * If you want to set reacetrack URL to something other than the default,
 * set it thru your system properties. For example:
 *
 *   java -DRACETRACK_URL=yourserver.eng.vmware.com com.your.app
 *
 * 2. (Deparecated, keeping it here for compatibility with some exising code)
 * If you already know your testCaseId, and can identify where you make your
 * first getInstance() call, then call getInstance(racetrackUrl, resultId)
 * first. Later on, always call RacetrackWebservice.getInstance() with no
 * arguments to get the Singleton instance of this class.
 * In this case, if you are not sure when you first getInstance() call is made,
 * always use getInstance(reacetrackUrl, resultId) to retrieve the Singleton.
 *
 * Test your additional changes with RacetrackWebserviceTest class
 */

public class RacetrackWebservice implements IRacetrack {

   private static final Logger _logger =
         LoggerFactory.getLogger(RacetrackWebservice.class);

   // Base URL of the Racetrack server.
   private String racetrackUrl = "https://racetrack.eng.vmware.com";
   private static final String RACETRACK_URL_SYS_PROPERTY = "RACETRACK_URL";
   private static final String HTTPS_PROTOCOL_PREFIX = "https://";

   // URL constants
   private static final String TEST_SET_BEGIN = "/TestSetBegin.php";
   private static final String TEST_SET_UPDATE = "/TestSetUpdate.php";
   private static final String TEST_SET_END = "/TestSetEnd.php";
   private static final String TEST_SET_DATA = "/TestSetData.php";
   private static final String TEST_CASE_BEGIN = "/TestCaseBegin.php";
   private static final String TEST_CASE_UPDATE = "/TestCaseUpdate.php";
   private static final String TEST_CASE_END = "/TestCaseEnd.php";
   private static final String COMMENT = "/TestCaseComment.php";
   private static final String VERIFICATION = "/TestCaseVerification.php";
   private static final String TEST_CASE_SCREENSHOT = "/TestCaseScreenshot.php";
   private static final String TEST_CASE_LOG = "/TestCaseLog.php";

   /**
    * Set this variable if you already know your resultId
    *
    * @deprecated always give testCaseId as an argument.
    */
   @Deprecated
   private NameValuePair resultId;

   /**
    * Initialise racetrack connection service.
    */
   public RacetrackWebservice(RacetrackConnectionSpec connectionSpec) {
      this.setRacetrackUrl(getValidatedUrl(connectionSpec.getTestLoggerURL()));
   }

   /**
    * Validates the passed in URL and attempts to convert it to a valid one if
    * validation fails.
    *
    * The method will throw an unchecked exception if both validations fail.
    *
    * @param connectionUrl - The URL to check.
    * @return
    */
   private String getValidatedUrl(String connectionUrl) {

      if (!isUrlValid(connectionUrl)) {
         _logger.info(
               "Racetrack url {} is not a valid URL. Attempting to add the {} prefix and validate again",
               connectionUrl,
               HTTPS_PROTOCOL_PREFIX
            );
         // Attempt to make the URL valid so as to preserve the behavior of
         // existing racetrack configurations.
         if (isUrlValid(HTTPS_PROTOCOL_PREFIX + connectionUrl)) {
            connectionUrl = HTTPS_PROTOCOL_PREFIX + connectionUrl;
         } else {
            // Throw an unchecked exception and do not proceed. The user needs
            // to know about this.
            throw new AssertionError(String.format(
                  "Racetrack url %s is not a valid URL", connectionUrl));
         }
      }

      return connectionUrl;
   }

   /**
    * Return whether the passed in urlAddress is a valid URL.
    *
    * @param urlAddress - The URL address to check.
    * @return - true if the address is a valid URL, false otherwise.
    */
   private boolean isUrlValid(final String urlAddress) {
      boolean validUrl;

      // The following line of code performs the validation, throwing an
      // exception on failure.
      try {
         URL url = new URL(urlAddress);
         url.toURI();
         validUrl = true;
      } catch (MalformedURLException | URISyntaxException e) {
         validUrl = false;
      }

      return validUrl;
   }

   /**
    * Create a TestSet in racetrack. Refer to
    * https://wiki.eng.vmware.com/RacetrackWebServices
    * for more details about the input arguments
    *
    * @param buildId
    * @param user
    * @param product
    * @param description
    * @param hostOs
    * @param serverBuildId
    * @param branch
    * @param branchType
    * @param testType
    * @param language
    * @return newly created TestSetID
    * @throws Exception
    */
   public String testSetBegin(String buildId, String user, String product,
         String description, String hostOs, String serverBuildId, String branch,
         String branchType, String buildType, String testType, String language)
         throws Exception {
      checkParamValidity("buildId", buildId);
      checkParamValidity("user", user);
      checkParamValidity("product", product);
      checkParamValidity("description", description);
      checkParamValidity("hostOs", hostOs);

      NameValuePair[] postData =
            new NameValuePair[] { new NameValuePair("BuildID", buildId),
                  new NameValuePair("User", user),
                  new NameValuePair("Product", product),
                  new NameValuePair("Description", description),
                  new NameValuePair("HostOS", hostOs),
                  new NameValuePair("ServerBuildID", serverBuildId),
                  new NameValuePair("Branch", branch),
                  new NameValuePair("BranchType", branchType),
                  new NameValuePair("BuildType", buildType),
                  new NameValuePair("TestType", testType),
                  new NameValuePair("Language", language) };

      return postRequest(
            getRacetrackUrl() + TEST_SET_BEGIN,
            purgePostData(postData)
         );
   }

   /**
    * Update a TestSet in racetrack. Refer to
    * https://wiki.eng.vmware.com/RacetrackWebServices
    * for more details about the input arguments
    *
    * @TODO: consider using a hash as input arg, so users
    * of this method only have to specify the params they
    * are changing. See ticket 373374 for more details
    * http://bugzilla.eng.vmware.com/show_bug.cgi?id=373374
    *
    * @param id ID of this test set to update
    * @param updatedValues hashtable of <key, value> pairs
    *        with test set property names and new values.
    *        Supported keys are case-sentitive and can be
    *        one of more of the following:
    *        "BuildID", "User", "Product", "Description",
    *        "HostOS", "ServerBuildID", "Branch", "BranchType",
    *        "BuildType", "TestType", and "Language"
    *
    * @throws TestException
    */
   public String testSetUpdate(String id, Hashtable<String, String> updatedValues)
         throws Exception {
      this.checkParamValidity("id", id);
      NameValuePair[] postData = new NameValuePair[updatedValues.size() + 1];
      int i = 0;
      postData[i++] = new NameValuePair("ID", id);
      Enumeration<String> keys = updatedValues.keys();
      while (keys.hasMoreElements()) {
         String key = keys.nextElement();
         postData[i++] = new NameValuePair(key, updatedValues.get(key));
      }
      return postRequest(this.getRacetrackUrl() + TEST_SET_UPDATE, postData);
   }

   /**
    * Ends a test set
    *
    * @param testSetId
    * @throws TestException
    */
   public void testSetEnd(String testSetId) throws Exception {
      this.checkParamValidity("testSetId", testSetId);
      NameValuePair[] postData =
            new NameValuePair[] { new NameValuePair("ID", testSetId) };
      postRequest(this.getRacetrackUrl() + TEST_SET_END, postData);
   }

   /**
    * Adds a name-value pair of data to a test set. Refer to the wiki
    * for more info:
    * https://wiki.eng.vmware.com/RacetrackWebServices#TestSetData
    *
    * @param testSetId
    * @param name
    * @param value
    * @return ResultSetDataId
    * @throws TestException
    */
   public String testSetData(String testSetId, String name, String value)
         throws Exception {
      this.checkParamValidity("testSetId", testSetId);
      this.checkParamValidity("Name", name);
      this.checkParamValidity("Value", value);
      NameValuePair[] data =
            { new NameValuePair("ResultSetID", testSetId),
                  new NameValuePair("Name", name), new NameValuePair("Value", value), };
      return postRequest(this.getRacetrackUrl() + TEST_SET_DATA, data);
   }

   /**
    *This method overloads the below method in order not to break existing clients
    *
    * @param testSetId
    * @param name
    * @param feature
    * @param description
    * @param machineName
    * @param tcmsId use it only if you are using testlink
    * @return newly created testCaseId
    * @throws TestException
    */
   public String testCaseBegin(String testSetId, String name, String feature,
         String description, String machineName, String tcmsId) throws Exception {
      return this.testCaseBegin(
            testSetId,
            name,
            feature,
            description,
            machineName,
            tcmsId,
            "EN");
   }

   /**
    * Starts a test case in a specific test set
    * Refer to the following wiki for more details:
    * https://wiki.eng.vmware.com/RacetrackWebServices#TestCaseBegin
    *
    * @param testSetId
    * @param name
    * @param feature
    * @param description
    * @param machineName
    * @param tcmsId use it only if you are using testlink
    * @param inputLanguage
    * @return newly created testCaseId
    * @throws TestException
    */
   public String testCaseBegin(String testSetId, String name, String feature,
         String description, String machineName, String tcmsId, String inputLanguage)
         throws Exception {
      this.checkParamValidity("testSetId", testSetId);
      this.checkParamValidity("name", name);
      this.checkParamValidity("feature", feature);
      NameValuePair[] postData =
            new NameValuePair[] { new NameValuePair("ResultSetID", testSetId),
                  new NameValuePair("Name", name),
                  new NameValuePair("Feature", feature),
                  new NameValuePair("Description", description),
                  new NameValuePair("MachineName", machineName),
                  new NameValuePair("TCMSID", tcmsId),
                  new NameValuePair("InputLanguage", inputLanguage) };
      return postRequest(
            this.getRacetrackUrl() + TEST_CASE_BEGIN,
            this.purgePostData(postData));
   }

   /**
    * Updates a test case
    *
    * Refer to the following wiki for more details:
    * https://wiki.eng.vmware.com/RacetrackWebServices#TestCaseBegin
    *
    * Only the id is required, all other parameters will be updated if they are not null.
    *
    * @param testCaseId The id of the test case to be updated
    * @param name Name of the test
    * @param feature Feature
    * @param description
    * @param machineName
    * @param tcmsId use it only if you are using testlink
    * @param inputLanguage
    * @return newly created testCaseId
    * @throws TestException
    */
   public String testCaseUpdate(String testCaseId, String name, String feature,
         String description, String machineName, String tcmsId, String inputLanguage,
         String guestOs) throws Exception {
      checkParamValidity("testCaseId", testCaseId);
      List<NameValuePair> pairList = new ArrayList<NameValuePair>();
      pairList.add(new NameValuePair("ID", testCaseId));
      if (name != null) {
         pairList.add(new NameValuePair("Name", name));
      }
      if (feature != null) {
         pairList.add(new NameValuePair("Feature", feature));
      }
      if (description != null) {
         pairList.add(new NameValuePair("Description", description));
      }
      if (machineName != null) {
         pairList.add(new NameValuePair("MachineName", machineName));
      }
      if (tcmsId != null) {
         pairList.add(new NameValuePair("TCMSID", tcmsId));
      }
      if (inputLanguage != null) {
         pairList.add(new NameValuePair("InputLanguage", inputLanguage));
      }
      if (guestOs != null) {
         pairList.add(new NameValuePair("GOS", guestOs));
      }
      NameValuePair[] postData = pairList.toArray(new NameValuePair[0]);

      return postRequest(getRacetrackUrl() + TEST_CASE_UPDATE, purgePostData(postData));
   }

   /**
    * Ends a test case
    *
    * @param testCaseId
    * @param result either "PASS" or "FAIL"
    * @throws TestException
    */
   public void testCaseEnd(String testCaseId, String result) throws Exception {
      this.checkParamValidity("testCaseId", testCaseId);
      this.checkParamValidity("result", result);
      NameValuePair[] postData =
            new NameValuePair[] { new NameValuePair("ID", testCaseId),
                  new NameValuePair("Result", result) };
      postRequest(this.getRacetrackUrl() + TEST_CASE_END, postData);
   }

   /**
    * Add verification info to racetrack for this test.
    *
    * @param description
    *           The description of the verification performed.
    * @param actual
    *           The actual result of the test
    * @param expected
    *           The expected result of the test
    * @param result
    *           The result of the verification true if it passed.
    *
    * @throws TestException
    *            thrown if anything goes wrong.
    */
   @Override
   public void testCaseVerification(String testCaseId, String description,
         String actual, String expected, boolean result) throws Exception {
      this.checkParamValidity("testCaseId", testCaseId);
      this.checkParamValidity("description", description);
      this.checkParamValidity("actual", actual);
      this.checkParamValidity("expected", expected);
      NameValuePair[] postData =
            new NameValuePair[] { new NameValuePair("ResultID", testCaseId),
                  new NameValuePair("Description", description),
                  new NameValuePair("Actual", actual),
                  new NameValuePair("Expected", expected),
                  new NameValuePair("Result", result ? "TRUE" : "FALSE") };
      postRequest(this.getRacetrackUrl() + VERIFICATION, postData);
   }

   /**
    * Add verification info to racetrack for this test.
    * Keeping this method for backward compatibility.
    * To use this method, you must initialize resultId
    * class variable
    *
    * @param description
    * @param actual
    * @param expected
    * @param result
    * @throws TestException
    *            thrown if anything goes wrong.
    * @deprecated  replaced by
    *  {@link #testCaseVerification(String,String,String,String,boolean)}
    */
   @Override
   @Deprecated
   public void testCaseVerification(String description, String actual, String expected,
         boolean result) throws Exception {
      if (this.getResultId() == null) {
         throw new Exception("Must specify testCaseId (ResultID)");
      }
      this.testCaseVerification(
            this.getResultId().getValue(),
            description,
            actual,
            expected,
            result);
   }

   /**
    * Log a comment in racetrack for this test.
    * Keeping this method for backward compatibility.
    * To use this method, you must initialize resultId
    * class variable
    *
    * @param comment a test case comment for this test.
    *
    * @throws TestException thrown if anything goes wrong.
    */
   @Override
   public void testCaseComment(String testCaseId, String comment) throws Exception {
      this.checkParamValidity("testCaseId", testCaseId);
      this.checkParamValidity("coment", comment);
      NameValuePair[] postData =
            new NameValuePair[] { new NameValuePair("ResultID", testCaseId),
                  new NameValuePair("Description", comment) };
      postRequest(this.getRacetrackUrl() + COMMENT, postData);
   }

   /**
    * Log a comment in racetrack for this test.
    *
    * @param comment a test case comment for this test.
    * @throws TestException thrown if anything goes wrong.
    * @deprecated replaced by {@link #testCaseComment(String,String)}
    */
   @Override
   @Deprecated
   public void testCaseComment(String comment) throws Exception {
      if (this.getResultId() == null) {
         throw new Exception("Must specify testCaseId (ResultID)");
      }
      this.testCaseComment(this.getResultId().getValue(), comment);
   }

   /**
    * Posts a screenshot to a testcase
    *
    * @param testCaseId
    * @param description
    * @param screenshot
    * @throws TestException
    */
   public void testCaseScreenshot(String testCaseId, String description,
         String screenshot) throws Exception {
      this.checkParamValidity("testCaseId", testCaseId);
      this.checkParamValidity("description", description);
      this.checkParamValidity("screenshot", screenshot);
      NameValuePair[] postData =
            new NameValuePair[] { new NameValuePair("ResultID", testCaseId),
                  new NameValuePair("Description", description) };
      postFileRequest(
            this.getRacetrackUrl() + TEST_CASE_SCREENSHOT,
            postData,
            "Screenshot",
            new File(screenshot));
   }

   /**
    * Posts a log file to a testcase
    *
    * Known issues: if log file is rempty, web service
    * will fail with "Bad Request" error.
    *
    * @param testCaseId
    * @param description
    * @param log
    * @throws TestException
    */
   public void testCaseLog(String testCaseId, String description, String log)
         throws Exception {
      this.checkParamValidity("testCaseId", testCaseId);
      this.checkParamValidity("description", description);
      this.checkParamValidity("log", log);
      NameValuePair[] postData =
            new NameValuePair[] { new NameValuePair("ResultID", testCaseId),
                  new NameValuePair("Description", description) };
      postFileRequest(this.getRacetrackUrl() + TEST_CASE_LOG, postData, "Log", new File(
            log));
   }

   /**
    * Post to Racetrack.
    *
    * @param url the url to post to.
    * @param data the request data
    * @return TestSetID or TestCaseID if any, "" otherwise
    *
    * @throws TestException thrown if anything goes wrong.
    */
   private String postRequest(String url, NameValuePair[] data) throws Exception {
      return postFileRequest(url, data, null, null);
   }

   /**
    * Post a multipart file upload request to Racetrack.
    *
    * Only one file per upload is supported, as limited by
    * current web service API
    *
    * @param url the url to post to.
    * @param data the request data
    * @return TestSetID or TestCaseID if any, "" otherwise
    *
    * @throws TestException thrown if anything goes wrong.
    */
   private String postFileRequest(String url, NameValuePair[] data, String uploadType,
         File fileToUpload) throws Exception {
      String result = "";
      HttpClient client = new HttpClient();
      client.getParams().setParameter("http.protocol.content-charset", "UTF-8");
      PostMethod post = new PostMethod(url);
      post.getParams().setParameter(
            HttpMethodParams.RETRY_HANDLER,
            new DefaultHttpMethodRetryHandler(3, false));
      if (uploadType == null || fileToUpload == null) {
         post.setRequestBody(data);
      } else {
         Part[] parts = new Part[data.length + 1];
         for (int i = 0; i < data.length; i++) {
            parts[i] = new StringPart(data[i].getName(), data[i].getValue());
         }
         try {
            parts[parts.length - 1] = new FilePart(uploadType, fileToUpload);
         } catch (FileNotFoundException ex) {
            throw new Exception(ex.getMessage());
         }
         post.setRequestEntity(new MultipartRequestEntity(parts, post.getParams()));
      }
      try {
         _logger.debug("Posting to URL: " + url);
         int status = client.executeMethod(post);

         if (HttpStatus.SC_OK != status) {
            throw new Exception("Got HTTP error: "
                  + HttpStatus.getStatusText(status));
         }
         String response = post.getResponseBodyAsString();
         result = response;
      } catch (HttpException e) {
         _logger.error(e.toString());
         throw new Exception("Caught HttpException");
      } catch (IOException e) {
         _logger.error(e.toString());
         throw new Exception("Caught IOException");
      }
      return result;
   }

   /**
    * Check validity of an input parameter.
    * If a param is required, then it may not be null or empty string,
    *
    * For now, we forbid passing in null as argument values.
    *
    * @param paramName
    * @param paramValue
    */
   public void checkParamValidity(String paramName, String paramValue)
         throws Exception {
      if (paramValue == null || paramValue.equals("")) {
         throw new Exception(paramName
               + " is a required field and may not be null or an empty string");
      }
   }

   /**
    * Remove NameValuePair with null or empty values from postData
    *
    * @param postData
    * @return new postData containing only valid NameValuePairs
    */
   public NameValuePair[] purgePostData(NameValuePair[] postData) {
      NameValuePair[] result = null;
      int count = 0;
      for (int i = 0; i < postData.length; i++) {
         if (postData[i].getValue() == null || postData[i].getValue().equals("")) {
            postData[i] = null;
         } else {
            count++;
         }
      }
      result = new NameValuePair[count];
      count = 0;
      for (int i = 0; i < postData.length; i++) {
         if (postData[i] != null) {
            result[count++] = postData[i];
         }
      }
      return result;
   }

   /**
    * @return the racetrackUrl
    */
   public String getRacetrackUrl() {
      return this.racetrackUrl;
   }

   /**
    * @param racetrackUrl the racetrackUrl to set
    */
   private void setRacetrackUrl(String racetrackUrl) {
      this.racetrackUrl = racetrackUrl;
   }

   /**
    * @return the resultId
    * @deprecated always give resultId as an argument
    */
   @Deprecated
   private NameValuePair getResultId() {
      return this.resultId;
   }
}
