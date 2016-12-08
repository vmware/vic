/* Copyright 2013 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.srv.common.srvapi;

import java.io.BufferedInputStream;
import java.io.BufferedOutputStream;
import java.io.IOException;
import java.io.InputStream;
import java.io.InputStreamReader;
import java.io.OutputStream;
import java.io.Reader;
import java.net.HttpURLConnection;
import java.net.URL;
import java.net.URLConnection;
import java.security.SecureRandom;
import java.security.Security;
import java.security.cert.CertificateException;
import java.security.cert.X509Certificate;

import javax.net.ssl.HostnameVerifier;
import javax.net.ssl.HttpsURLConnection;
import javax.net.ssl.SSLContext;
import javax.net.ssl.SSLSession;
import javax.net.ssl.TrustManager;
import javax.net.ssl.X509TrustManager;

import org.apache.commons.codec.binary.Base64;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import com.sun.net.ssl.internal.ssl.Provider;

/**
 * Takes care of transferring files over HTTP/S or FTP.
 */
public class HttpConnector {

   private static final Logger _logger = LoggerFactory.getLogger(HttpConnector.class);

   static {
      try {
         Security.addProvider(new Provider());

         //Create a trust manager that does not validate certificate chains.
         TrustManager[] trustAllCerts = new TrustManager[] {
               new X509TrustManager() {
                  @Override
                  public X509Certificate[] getAcceptedIssuers() {
                     return null;
                  }

                  @Override
                  public void checkServerTrusted(X509Certificate[] cert, String authType)
                        throws CertificateException {
                  }

                  @Override
                  public void checkClientTrusted(X509Certificate[] cert, String authType)
                        throws CertificateException {
                  }
               }
         };

         SSLContext sc = SSLContext.getInstance("SSL");
         sc.init(null, trustAllCerts, new SecureRandom());
         HttpsURLConnection.setDefaultSSLSocketFactory(sc.getSocketFactory());

         final HostnameVerifier hostNameVerifier =
               new HostnameVerifier() {
                  @Override
                  public boolean verify(final String host, final SSLSession session) {
                     _logger.debug("HostnameVerifier: " + host + " - " + session);
                     return true;
                  }
               };
         HttpsURLConnection.setDefaultHostnameVerifier(hostNameVerifier);
      } catch (Exception exception) {
         _logger.error("Unable to initialize security trust manager/hostname verifier");
         throw new RuntimeException(exception);
      }
   }

   /**
    * Upload file from one URL to another.
    *
    * @param srcURL        source URL from where file will be copied
    * @param destURL       destination URL where the file will be copied
    * @return              if the operation was successful
    * @throws Exception    if there were any problems during connections or transfer
    */
   public static boolean uploadToServer(URL srcURL, URL destURL) throws Exception {
      boolean uploaded = false;
      InputStream in = null;
      OutputStream out = null;
      HttpURLConnection destConn = null;
      HttpURLConnection srcConn = null;
      URLConnection ftpurlConn = null;

      try {
         destConn = getUploadPostConn(destURL);
         destConn.setChunkedStreamingMode(1024 * 10);
         destConn.connect();
         out = destConn.getOutputStream();
         if (srcURL.getProtocol().equalsIgnoreCase("http")
               || srcURL.getProtocol().equalsIgnoreCase("https")) {
            srcConn = getDownloadConn(srcURL);
            in = srcConn.getInputStream();
         } else if (srcURL.getProtocol().equalsIgnoreCase("ftp")) {
            ftpurlConn = srcURL.openConnection();
            in = ftpurlConn.getInputStream();
         } else {
            throw new IllegalArgumentException(
                  "The source URL protocol is not supported"
               );
         }
         uploaded = copy(srcConn, destConn, in, out);
      } finally {
         if (in != null) {
            in.close();
         }
         if (out != null) {
            out.close();
         }
         if (srcConn != null) {
            srcConn.disconnect();
         }
         if (destConn != null) {
            destConn.disconnect();
         }
      }
      return uploaded;
   }

   /**
    * Retrieves a file from the passed URL.
    *
    * @param fileUrl        the URL from where the file should be retrieved
    * @return              the contents of the file itself
    * @throws IOException  if something wrong happens during reading
    */
   public static String getFileFromUrl(String fileUrl) throws IOException {
      Reader reader = new InputStreamReader(new URL(fileUrl).openStream());
      StringBuffer buffer = new StringBuffer();
      while (true) {
         int i = reader.read();
         if (i == -1) {
            break;
         }
         buffer.append((char) i);
      }

      return buffer.toString();
   }

   // ---------------------------------------------------------------------------
   // Private methods

   private static HttpURLConnection getUploadPostConn(URL url) throws IOException {
      final HttpURLConnection httpConn = getConnection(url);
      httpConn.setDoOutput(true);
      httpConn.setRequestProperty("Content-Type", "application/x-vnd.vmware-streamVmdk");
      httpConn.setRequestProperty("Expect", "    -continue");
      httpConn.setRequestProperty("Overwrite", "t");
      httpConn.setRequestMethod("POST");
      return httpConn;
   }

   private static HttpURLConnection getConnection(final URL url) throws IOException {
      HttpURLConnection httpConn = null;
      httpConn = (HttpURLConnection) url.openConnection();
      httpConn.setDoInput(true);
      httpConn.setAllowUserInteraction(true);
      httpConn.setReadTimeout(60000);
      httpConn.setRequestProperty(
            "Authorization",
            " Basic " + Base64.encodeBase64String(":".getBytes())
         );
      return httpConn;
   }

   private static HttpURLConnection getDownloadConn(final URL url , String ...id)
         throws IOException {
      final HttpURLConnection httpConn = getConnection(url);
      if (id != null && id.length > 0) {
         try {
            httpConn.setRequestProperty("Cookie", id[0]);
         } catch (java.lang.IllegalStateException ise) {
            _logger.error("failed to set cookie" +ise.toString());
         }
      }
      httpConn.connect();
      _logger.info("getDownloadConn() - HTTP Method: "
               + httpConn.getRequestMethod() + "  Response: "
               + httpConn.getResponseCode());

      int responseCode = httpConn.getResponseCode();
      if (responseCode >= 200 && responseCode < 300) {
         _logger.info("Successfully got download connection.");
      } else {
         throw new IOException("Response code: " + responseCode);
      }

      return httpConn;
   }

   private static boolean copy(HttpURLConnection sourceConnection,
         HttpURLConnection destinationConnection, InputStream in, OutputStream out)
         throws IOException {
      boolean copied = false;
      if (copy(in, out)) {
         final int respCode = destinationConnection.getResponseCode();
         final String responseMsg = destinationConnection.getResponseMessage();
         _logger.info("Response Code: " + respCode + ":" + responseMsg);
         if (respCode == HttpURLConnection.HTTP_OK) {
            _logger.info("Uploaded file might have been replaced on the server");
            copied = true;
         } else if (respCode == HttpURLConnection.HTTP_CREATED) {
            _logger.info("Uploaded file created on server");
            copied = true;
         } else {
            throw new IOException("Failed to upload file - response code: " + respCode
                  + "  Message" + responseMsg);
         }
      } else {
         _logger.error("Failed to copy");
      }
      return copied;
   }

   private static boolean copy(final InputStream in, final OutputStream out) {
      final long startTime = System.currentTimeMillis();
      long endTime = startTime;
      boolean result = false;
      final byte[] buf = new byte[1024 * 10];
      long len = 0;
      _logger.info("Copy - Started copying.... Using buffer size 1024 * 10:");
      final BufferedInputStream bin = new BufferedInputStream(in);
      final BufferedOutputStream bout =
            new BufferedOutputStream(out, 1024 * 10);
      try {
         int read = 0;
         int print = 0;
         while ((read = bin.read(buf)) > 0) {
            bout.write(buf, 0, read);
            bout.flush();
            len += read;
            print++;
            if (print >= 50) {
               _logger.info("> " + len);
               print = 0;
            }
         }
         endTime = System.currentTimeMillis();
         _logger.info(len + " of data transferred to destination. Time taken : "
               + ((float) (endTime - startTime) / (0)) + " seconds");
         result = true;
      } catch (Exception e) {
         _logger.error("Unable to copy: " + e.getMessage());
      }
      return result;
   }
}
