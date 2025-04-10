package com.example.totp;

import javax.crypto.Mac;
import javax.crypto.spec.SecretKeySpec;

import java.nio.ByteBuffer;
import java.nio.ByteOrder;
import java.time.LocalDateTime;
import java.time.ZoneOffset;
import java.time.ZonedDateTime;
import java.time.format.DateTimeFormatter;
import java.util.Arrays;
import java.util.Base64;
import android.util.Log;
import android.content.Context;

public class Totp {
    private static final DateTimeFormatter DateTimeFormat = DateTimeFormatter.ofPattern("yyyy-MM-dd HH:mm:ss");
    private static final String ClientSecretKey = "clientsecretkey";
    private int otp = 0;
    private String serverSecret = BuildConfig.totpSecret;
    private String clientSecret = "";
    private int interval = 5;
    private int digit = 6;

    private final KeyManager keyManager;
    private State state;

    public Totp(KeyManager keyManager) {
        this.keyManager = keyManager;
        this.state = State.Initial;
    }

    private static enum State {
        Initial,
        WaitForLogin,
        Idle
    }

    public OtpStatusCode updateStatus() {
        OtpStatusCode result;
        if (state == State.Initial) {
            state = State.WaitForLogin;
            String secret = keyManager.get(ClientSecretKey);
            if (secret == null) {
                state = State.WaitForLogin;
                result = OtpStatusCode.LoginRequired;
            }else {
                state = State.Idle;
                clientSecret = secret;
                result = OtpStatusCode.Ok;
            }
        }else if (state == State.WaitForLogin) {
            result = OtpStatusCode.LoginRequired;
        }else {
            result = OtpStatusCode.Ok;
        }
        return result;
    }

    public OtpResult generateOtp() {
        OtpStatusCode result = updateStatus();
        if (state != State.Idle) {
            return new OtpResult(result, 0);
        }

        int max = (int) Math.pow(10, digit);
        String timeString = getRoundedTimeString();
        String data = timeString + clientSecret;
        byte[] hashed;
        try {
            hashed = computeHash(data.getBytes());
        }catch (Exception e) {
            Log.d("Failed to compute hash: ", e.toString());
            return new OtpResult(OtpStatusCode.InternalError, 0);
        }
        int truncated = truncate(hashed);
        int rs = truncated % max;
        // Log.d("lala: ", data + " - " + Integer.toString(rs));
        return new OtpResult(OtpStatusCode.Ok, rs);
    }

    private String getRoundedTimeString() {
        //todo: use utc time
        LocalDateTime now = ZonedDateTime.now(ZoneOffset.UTC).toLocalDateTime();
        int nowRounded = now.getSecond() - now.getSecond() % interval;
        
        String timeString = now.withSecond(nowRounded).format(DateTimeFormat);
        return timeString;
    }

    private byte[] computeHash(byte[] data) throws Exception {
        SecretKeySpec keySpec = new SecretKeySpec(serverSecret.getBytes(), "HmacSHA256");
        Mac mac = Mac.getInstance("HmacSHA256");
        mac.init(keySpec);
        return mac.doFinal(data);
    }

    private int truncate(byte[] data) {
        byte[] truncated = Arrays.copyOfRange(data, 0, 4);
        truncated[3] = 0;

        ByteBuffer buffer = ByteBuffer.wrap(truncated);
        buffer.order(ByteOrder.LITTLE_ENDIAN);
        return buffer.getInt();
    }

    public void setSeed(String seed) {
        keyManager.set(ClientSecretKey, seed);
        clientSecret = seed;
        state = State.Idle;
    }
}
