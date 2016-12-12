package com.vmware.client.automation.util.ssh;

import java.io.IOException;

import ch.ethz.ssh2.Connection;
import ch.ethz.ssh2.InteractiveCallback;

/**
 * <code>SSHConnection</code> is ... <br>
 * TODO 1. Describe the class in one sentence <br>
 * TODO 2. Describe the purpose of this class in detail <br>
 * TODO 3. (optional) Describe any concurrency considerations <br>
 * TODO 4. (optional) Give examples of usage if needed <br>
 * 
 * @since <Product revision>
 * @version <Implementation version of this type>
 * @author dproynov
 */
public class SshConnectionManager implements ConnectionManager {

   private final String host;
   private final String username;
   private final String password;

   private SshConnection defaultConnection;

   /**
    * @param host
    * @param userName
    * @param password
    * 
    */
   public SshConnectionManager(String host, String userName, String password) {
      this.host = host;
      this.username = userName;
      this.password = password;
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
    * @param host
    * @param userName
    * @param password
    * @return
    * @throws IOException
    */
   private SshConnection connectInternal() throws IOException {
      boolean authenticated = false;
      Connection connInternal = new Connection(host);
      connInternal.connect();

      if (connInternal.isAuthMethodAvailable(username, "password")) {
         authenticated = connInternal.authenticateWithPassword(username, password);
      } else if (connInternal.isAuthMethodAvailable(username, "keyboard-interactive")) {
         authenticated = connInternal.authenticateWithKeyboardInteractive(username,
               new InteractiveCallback() {
                  @Override
                  public String[] replyToChallenge(String name, String instruction,
                        int numPrompts, String[] prompt, boolean[] echo)
                              throws Exception {
                     String[] responses = new String[numPrompts];
                     for (int i = 0; i < numPrompts; i++) {
                        responses[i] = password;
                     }
                     return responses;
                  }
               });
      }

      if (!authenticated) {
         throw new IOException("Authentication failed");
      }
      SshConnection conn = new SshConnection(connInternal);
      return conn;
   }

   @Override
   public void close() {
      if (defaultConnection != null) {
         defaultConnection.close();
      }
   }

   /**
    * @return the defaultConnection
    */
   public SshConnection getDefaultConnection() {
      if (defaultConnection == null) {
         try {
            setDefaultConnection(connectInternal());
         } catch (IOException e) {
            e.printStackTrace();
         }
      }
      return defaultConnection;
   }

   /**
    * @param defaultConnection
    *           the defaultConnection to set
    */
   public void setDefaultConnection(SshConnection defaultConnection) {
      this.defaultConnection = defaultConnection;
   }

}
