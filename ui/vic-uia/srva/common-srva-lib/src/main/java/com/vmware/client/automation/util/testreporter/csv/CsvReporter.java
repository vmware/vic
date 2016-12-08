/* Copyright 2014 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.client.automation.util.testreporter.csv;

import java.io.BufferedWriter;
import java.io.File;
import java.io.FileOutputStream;
import java.io.IOException;
import java.io.OutputStreamWriter;
import java.io.Writer;
import java.text.DateFormat;
import java.text.SimpleDateFormat;
import java.util.Date;
import java.util.List;
import java.util.Map;
import java.util.Set;

import org.testng.IReporter;
import org.testng.IResultMap;
import org.testng.ISuite;
import org.testng.ISuiteResult;
import org.testng.ITestContext;
import org.testng.ITestNGMethod;
import org.testng.ITestResult;
import org.testng.xml.XmlSuite;

import com.vmware.client.automation.workflow.BaseTest.TestID;

/**
 * Custom test ng results reporter.
 * Reports results in comma separated values (CSV) file format.
 */
public class CsvReporter implements IReporter {
   private static final String FILE_NAME = "results_%s_%s.csv";
   private static final String C = ",";
   private static final String NL = System.lineSeparator();

   private final StringBuilder csvBuilder =
         new StringBuilder("ID,TestName,Status,Groups" + NL);

   private static final String PASSED = "PASSED";
   private static final String FAILED = "FAILED";
   private static final String SKIPPED = "SKIPPED";
   private static final String UNKNOWN = "UNKNOWN";

   @Override
   public void generateReport(List<XmlSuite> xmlSuites,
                              List<ISuite> suites,
                              String outputDirectory) {
      for (ISuite suite : suites) {
         Map<String, ISuiteResult> resultMap = suite.getResults();

         for (ISuiteResult result : resultMap.values()) {
            ITestContext tc = result.getTestContext();

            append(tc.getPassedTests());
            append(tc.getFailedTests());
            append(tc.getSkippedTests());
         }

         Writer writer = null;
         DateFormat dateFormat = new SimpleDateFormat("yyyy-MM-dd_HH-mm-ss");
         File file = new File(outputDirectory + File.separator +
                              String.format(FILE_NAME,
                                            suite.getName(),
                                            dateFormat.format(new Date())));
         try {
            writer = new BufferedWriter(
                  new OutputStreamWriter(new FileOutputStream(file), "utf-8"));
            writer.write(csvBuilder.toString());
         } catch (IOException e) {
            e.printStackTrace();
         } finally {
            try {
               writer.close();
            } catch (Exception ex) {
               // Nothing to do.
            }
         }

      }
   }

   private void append(IResultMap results) {
      for (ITestNGMethod m : results.getAllMethods()) {
         Set<ITestResult> methodResults = results.getResults(m);
         for (ITestResult itr : methodResults) {
            append(m.getConstructorOrMethod().getMethod().getAnnotation(
                  TestID.class).id()[0] + C);
            append(m.getMethodName() + C);
            append(getStatusString(itr.getStatus()) + C);
            append(m.getGroups());
            append(NL);
         }
      }
   }

   private void append(String txt) {
      csvBuilder.append(txt);
   }

   private void append(String[] txt) {
      for (int i = 0; i < txt.length; i++) {
         csvBuilder.append(txt[i]);
         if (i < txt.length - 1) {
            csvBuilder.append(";");
         }
      }
   }

   private String getStatusString(int statusCode) {
      switch (statusCode) {
         case ITestResult.SUCCESS:
            return PASSED;
         case ITestResult.FAILURE:
            return FAILED;
         case ITestResult.SKIP:
            return SKIPPED;
         default:
            return UNKNOWN + ":" + statusCode;
      }
   }

}
