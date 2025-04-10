package com.example.totp;

public enum OtpStatusCode {
    Ok(0),
    LoginRequired(1),
    InternalError(2),
    ServerError(3),
    Unknow(4);

    private int value;

    private OtpStatusCode(int value) {
        this.value = value;
    }

    public int getValue() {
        return this.value;
    }
}
