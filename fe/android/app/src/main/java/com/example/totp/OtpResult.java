package com.example.totp;

public class OtpResult {
    private final OtpStatusCode statusCode;
    private final int otp;

    public OtpResult(OtpStatusCode statusCode, int otp) {
        this.statusCode = statusCode;
        this.otp = otp;
    }

    public OtpStatusCode getStatusCode() {
        return this.statusCode;
    }

    public int getOtp() {
        return this.otp;
    }
}
