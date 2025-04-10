package com.example.totp;

import java.io.*;
import java.net.HttpURLConnection;
import java.net.URL;
import javax.net.ssl.*;
import java.security.SecureRandom;
import java.security.cert.X509Certificate;
import java.nio.charset.StandardCharsets;
import java.util.concurrent.ExecutorService;
import java.util.concurrent.Executors;
import java.util.concurrent.Future;

import org.json.JSONObject;
import android.util.Log;

public class LoginHandler {
    //todo: load from env file
    private static final String server = "https://172.19.175.104:8080";
    private static final String ClientTokenKey = "clienttoken";

    private final KeyManager keyManager;
    private final Totp totp;
    private final ExecutorService executorService;
    private static boolean disableSslVerification = false;

    public LoginHandler(KeyManager keyManager, Totp totp) {
        this.keyManager = keyManager;
        this.totp = totp;
        this.executorService = Executors.newSingleThreadExecutor();
        disableSslVerification();
    }

    private static void disableSslVerification(){
        if (disableSslVerification) {
            return;
        }
        disableSslVerification = true;
        // Create a trust manager that does not validate certificate chains
        try {
            TrustManager[] trustAllCerts = new TrustManager[]{
                new X509TrustManager() {
                    public X509Certificate[] getAcceptedIssuers() { return new X509Certificate[0]; }
                    public void checkClientTrusted(X509Certificate[] certs, String authType) {}
                    public void checkServerTrusted(X509Certificate[] certs, String authType) {}
                }
            };
    
            // Install the all-trusting trust manager
            SSLContext sc = SSLContext.getInstance("TLS");
            sc.init(null, trustAllCerts, new SecureRandom());
            HttpsURLConnection.setDefaultSSLSocketFactory(sc.getSocketFactory());
    
            // Create all-trusting host name verifier
            HostnameVerifier allHostsValid = (hostname, session) -> true;
    
            // Install the all-trusting host verifier
            HttpsURLConnection.setDefaultHostnameVerifier(allHostsValid);
        }catch (Exception e) {
            Log.d("disableSslVerification(): Failed to disable ssl: ", e.toString());
            disableSslVerification = false;
        }
        
    }

    public OtpStatusCode login(String username, String password) {
        Future<OtpStatusCode> result =  executorService.submit(
            () -> {
                return _login(username, password);
            }
        );
        OtpStatusCode ret;
        try {
            ret = result.get();
        }catch (Exception e) {
            return OtpStatusCode.InternalError;
        }
        if (ret != OtpStatusCode.Ok) {
            return ret;
        }

        result =  executorService.submit(
            () -> {
                return getSeed();
            }
        );

        try {
            ret = result.get();
        }catch (Exception e) {
            return OtpStatusCode.InternalError;
        }

        return ret;
    }

    private OtpStatusCode _login(String username, String password) {
        try {
            disableSslVerification();
            URL url = new URL(String.format("%s/login", server));
            HttpURLConnection connection = (HttpURLConnection) url.openConnection();
            connection.setRequestMethod("POST");
            connection.setRequestProperty("Content-Type", "application/json");
            connection.setDoOutput(true);
            JSONObject input = new JSONObject();
            input.put("username", username);
            input.put("password", password);
            OutputStream os = connection.getOutputStream();
            byte[] inputByte = input.toString().getBytes(StandardCharsets.UTF_8);
            os.write(inputByte, 0, inputByte.length);
    
            int responseCode = connection.getResponseCode();
            if (responseCode == HttpURLConnection.HTTP_OK) {
                BufferedReader in = new BufferedReader(new InputStreamReader(connection.getInputStream()));
                String inputLine;
                StringBuilder response = new StringBuilder();
                while((inputLine = in.readLine()) != null) {
                    response.append(inputLine);
                }
                in.close();
                JSONObject respJson = new JSONObject(response.toString());
                String access_token = respJson.getString("access_token");
                if (access_token == null) {
                    Log.d("_login(): cannot get token ??", "");
                    return OtpStatusCode.ServerError;
                }
                Log.d("login(): access token: ", access_token);
                keyManager.set(ClientTokenKey, access_token);
                return OtpStatusCode.Ok;
            }
        }catch (Exception e) {
            Log.d("login(): error: ", e.toString());
            return OtpStatusCode.InternalError;
        }

        return OtpStatusCode.LoginRequired;
    }

    private OtpStatusCode getSeed() {
        try {
            disableSslVerification();
            String token = keyManager.get(ClientTokenKey);
            if (token == null) {
                Log.d("getSeed(): cannot get token ??? ", "");
                return OtpStatusCode.LoginRequired;
            }
            URL url = new URL(String.format("%s/seed", server));
            HttpURLConnection connection = (HttpURLConnection) url.openConnection();
            connection.setRequestMethod("GET");
            connection.setRequestProperty("Accept", "application/json");
            connection.setRequestProperty("Authorization", "Bearer " + token);

            int responseCode = connection.getResponseCode();
            if (responseCode != HttpURLConnection.HTTP_OK) {
                Log.d("getSeed(): cannot get seed status code: ", Integer.toString(responseCode));
                return OtpStatusCode.LoginRequired;
            }
            BufferedReader in = new BufferedReader(new InputStreamReader(connection.getInputStream()));
            String inputLine;
            StringBuilder response = new StringBuilder();
            while((inputLine = in.readLine()) != null) {
                response.append(inputLine);
            }
            in.close();
            JSONObject respJson = new JSONObject(response.toString());
            String seed = respJson.getString("seed");
            if (seed == null) {
                Log.d("getSeed(): seed is null ??", "");
                return OtpStatusCode.LoginRequired;
            }
            totp.setSeed(seed);
        }catch (Exception e) {
            Log.d("login(): error: ", e.toString());
            return OtpStatusCode.InternalError;
        }
        return OtpStatusCode.Ok;
    }
}

