package com.vmware.client.automation.util.ssh;

import java.io.BufferedReader;
import java.io.IOException;
import java.io.InputStream;
import java.io.InputStreamReader;
import java.io.Writer;
import java.util.concurrent.CountDownLatch;
import java.util.concurrent.TimeUnit;

import ch.ethz.ssh2.Connection;
import ch.ethz.ssh2.Session;
import ch.ethz.ssh2.StreamGobbler;

/**
 * <code>SshConnection</code> is ... <br>
 * TODO 1. Describe the class in one sentence <br>
 * TODO 2. Describe the purpose of this class in detail <br>
 * TODO 3. (optional) Describe any concurrency considerations <br>
 * TODO 4. (optional) Give examples of usage if needed <br>
 *
 * @since <Product revision>
 * @version <Implementation version of this type>
 * @author dproynov
 */
public class SshConnection {

   private final ch.ethz.ssh2.Connection connection;

   /**
    * @param connection
    */
   public SshConnection(Connection connection) {
      this.connection = connection;
   }

   /**
    * Executes given ssh command and writes its output to the given writers.
    *
    * @param command - String representing the remote command
    * @param outWriter - Writer to write the standard out
    * @param errWriter - Writer to write the standard err
    * @param timeoutInSeconds - timeout in seconds
    * @throws IOException - If there is a problem with the command or the
    *                   connection or the timeout has been reached
    */
   public void executeSshCommand(String command, Writer outWriter, Writer errWriter,
         int timeoutInSeconds) throws IOException {
      // Only one command can be executed per ssh session -
      // restriction of ssh
      Session sshSession = connection.openSession();
      InputStream stdout = new StreamGobbler(sshSession.getStdout());
      InputStream stderr = new StreamGobbler(sshSession.getStderr());

      BufferedReader stdoutReader = new BufferedReader(new InputStreamReader(stdout));
      BufferedReader stderrReader = new BufferedReader(new InputStreamReader(stderr));

      CountDownLatch latch = new CountDownLatch(2);

      // execute command
      sshSession.execCommand(command);

      Thread stdoutConsumer =
            new Thread(new SshOutputConsumerRunnable(stdoutReader, outWriter, latch));
      Thread stderrConsumer =
            new Thread(new SshOutputConsumerRunnable(stderrReader, errWriter, latch));

      stdoutConsumer.start();
      stderrConsumer.start();

      boolean completedInTime = true;
      try {
         //blocking method
         completedInTime = latch.await(timeoutInSeconds, TimeUnit.SECONDS);
      } catch (InterruptedException e) {
         e.printStackTrace();
      } finally {
         // Clean up this session
         stdoutReader.close();
         stderrReader.close();
         stdout.close();
         stderr.close();
         sshSession.close();

         if (!completedInTime) {
            throw new IOException("The command: " + command
                  + " did not complete in the " + "given timeout (" + timeoutInSeconds
                  + " seconds)");
         }
      }
   }

   /**
    * Executes given ssh command and writes its output to the given writers.
    *
    * </p>
    * <b>WARNING! This method times out in 70 minutes by default. </b>
    * If you want to configure the timeout use the overloaded variant.
    *
    * @param command - String representing the remote command
    * @param outWriter - Writer to write the standard out
    * @param errWriter - Writer to write the standard err
    * @throws IOException - If there is a problem with the command or the
    *                   connection or the timeout has been reached
    */
   public void executeSshCommand(String command, Writer outWriter, Writer errWriter)
         throws IOException {
      executeSshCommand(command, outWriter, errWriter, 60 * 70);
   }

   /**
    *
    * TODO 1. Describe the method in one sentence <br>
    * TODO 2. Describe the method's purpose in detail <br>
    * TODO 3. Describe any requirements/pre-conditions, modifications, and
    * effects/post-conditions <br>
    * TODO 4. (optional) Describe any concurrency considerations <br>
    * TODO 5. (optional) Give examples of usage if needed <br>
    *
    */
   public void close() {
      connection.close();
   }

   private class SshOutputConsumerRunnable implements Runnable {

      private final BufferedReader input;
      private final Writer output;
      private final CountDownLatch latch;

      /**
       * @param input
       * @param output
       * @param latch
       *
       */
      public SshOutputConsumerRunnable(BufferedReader input, Writer output,
            CountDownLatch latch) {
         this.input = input;
         this.output = output;
         this.latch = latch;
      }

      @Override
      public void run() {
         try {
            while (true) {
               String outLine = input.readLine();
               if (outLine == null) {
                  latch.countDown();
                  break;
               }
               if (outLine != null && output != null) {
                  output.write(outLine + "\n");
                  output.flush();
               }
            }
         } catch (IOException e) {
            e.printStackTrace();
         }
      }

   }

}
