package com.example.totp;

import android.os.Bundle;
import io.flutter.embedding.android.FlutterActivity;
import io.flutter.plugin.common.MethodChannel;
import android.util.Log;
import java.util.HashMap;
import java.util.Map;


public class MainActivity extends FlutterActivity {
    private static String CHANNEL = "com.example.totp/channel";

    private Totp totp;
    private LoginHandler loginHandler;

    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        KeyManager keyManager = new KeyManager(this);
        totp = new Totp(keyManager);
        loginHandler = new LoginHandler(keyManager, totp);

        new MethodChannel(getFlutterEngine().getDartExecutor(), CHANNEL).setMethodCallHandler(
            (call, result) -> {
                if (call.method.equals("generateOtp")) {
                    OtpResult otpResult = generateOtp();

                    Map<String, Object> resultMap = new HashMap<>();
                    resultMap.put("result", otpResult.getStatusCode().getValue());
                    resultMap.put("otp", otpResult.getOtp());
    
                    result.success(resultMap);
                }else if (call.method.equals("updateStatus")) {
                    result.success(updateStatus());
                }else if (call.method.equals("login")) {
                    String username = call.argument("username");
                    String password = call.argument("password");
                    result.success(login(username, password));
                }
                else {
                    result.notImplemented();
                }
            }
        );
    }

    private OtpResult generateOtp() {
        return totp.generateOtp();
    }

    private int updateStatus() {
        return totp.updateStatus().getValue();
    }

    private int login(String username, String password) {
        return loginHandler.login(username, password).getValue();
    }
}