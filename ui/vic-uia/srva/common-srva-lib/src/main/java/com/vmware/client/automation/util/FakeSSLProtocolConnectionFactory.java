/* Copyright 2013 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.client.automation.util;

import java.security.cert.CertificateException;
import java.security.KeyManagementException;
import java.security.KeyStoreException;
import java.security.NoSuchAlgorithmException;
import java.security.cert.X509Certificate;

import javax.net.ssl.SSLContext;
import javax.net.ssl.TrustManager;
import javax.net.ssl.X509TrustManager;

import org.apache.http.config.Registry;
import org.apache.http.config.RegistryBuilder;
import org.apache.http.conn.socket.ConnectionSocketFactory;
import org.apache.http.conn.ssl.SSLConnectionSocketFactory;
import org.apache.http.conn.ssl.TrustSelfSignedStrategy;
import org.apache.http.conn.ssl.X509HostnameVerifier;
import org.apache.http.impl.client.CloseableHttpClient;
import org.apache.http.impl.client.HttpClientBuilder;
import org.apache.http.impl.client.HttpClients;
import org.apache.http.impl.conn.PoolingHttpClientConnectionManager;

import com.vmware.vapi.internal.protocol.client.rpc.CorrelatingClient;
import com.vmware.vapi.internal.protocol.client.rpc.http.HttpClient;
import com.vmware.vapi.protocol.ClientConfiguration;
import com.vmware.vapi.protocol.HttpConfiguration;
import com.vmware.vapi.protocol.JsonProtocolConnectionFactory;

/**
 * This class is a fake SSL ProtocolConnectionFactory. Here is configured a
 * HttpClient with a SSL context using trust self signed strategy and with a
 * ALLOW_ALL_HOSTNAME_VERIFIER. Using such HttpClient the UI Automation is
 * enabled to do connections to the https VAPI endpoint without the need to
 * download the certificate from the cloud VM. The only possible mechanism to
 * provide such a pre-configured HttpClient to the vapi is by implementing
 * ProtocolConnectionFactory and overwrite the createHttpTransport method.
 */
public class FakeSSLProtocolConnectionFactory extends JsonProtocolConnectionFactory {

   private final int MAX_TOTAL_CONNECTIONS = 2000;
   private final int MAX_CONN_PER_ROUTE = 600;
   private final String HTTPS_PROTOCOL = "https";

   private final org.apache.http.impl.client.CloseableHttpClient _client;

   public FakeSSLProtocolConnectionFactory() throws KeyManagementException,
   NoSuchAlgorithmException, KeyStoreException {
      this._client = buildHttpClient();
   }

   @Override
   protected CorrelatingClient createHttpTransport(String uri,
         ClientConfiguration clientConfig, HttpConfiguration httpConfig) {
      return new HttpClient(uri, this._client);
   }

   // ---------------------------------------------------------------------------
   // Private methods

   /**
    * Create a configured fake SSL HttpClient.
    *
    * @return
    * @throws KeyManagementException
    * @throws NoSuchAlgorithmException
    * @throws KeyStoreException
    */
   private org.apache.http.impl.client.CloseableHttpClient buildHttpClient()
         throws KeyManagementException, NoSuchAlgorithmException,
         KeyStoreException {

      Registry<ConnectionSocketFactory> schemeRegistry = buildAllowAllHostRegistry();

      PoolingHttpClientConnectionManager cm = new PoolingHttpClientConnectionManager(
            schemeRegistry);
      cm.setMaxTotal(MAX_TOTAL_CONNECTIONS);
      cm.setDefaultMaxPerRoute(MAX_CONN_PER_ROUTE);

      HttpClientBuilder builder = HttpClients.custom();
      builder.setConnectionManager(cm);
      return builder.build();
   }

   /**
    * Build a registry of self signed trust strategy SSL context.
    *
    * @return
    * @throws KeyManagementException
    * @throws NoSuchAlgorithmException
    * @throws KeyStoreException
    */
   private Registry<ConnectionSocketFactory> buildAllowAllHostRegistry()
         throws KeyManagementException, NoSuchAlgorithmException,
         KeyStoreException {
      RegistryBuilder<ConnectionSocketFactory> builder = RegistryBuilder
            .<ConnectionSocketFactory> create();

      SSLContext sslContext = SSLContext.getInstance("TLSv1.2");

      X509TrustManager allTrustingManager = new X509TrustManager() {
         public void checkClientTrusted(X509Certificate[] xcs,
               String string) throws CertificateException {
         }
      
         public void checkServerTrusted(X509Certificate[] xcs,
               String string) throws CertificateException {
         }
      
         public X509Certificate[] getAcceptedIssuers() {
            return null;
         }
      };

      sslContext.init(null, new TrustManager[] { allTrustingManager }, null);

      X509HostnameVerifier hostVerifier = SSLConnectionSocketFactory.ALLOW_ALL_HOSTNAME_VERIFIER;
      SSLConnectionSocketFactory sslcsf = new SSLConnectionSocketFactory(
            sslContext, hostVerifier);
      builder.register(HTTPS_PROTOCOL, sslcsf);
      return builder.build();
   }
}
