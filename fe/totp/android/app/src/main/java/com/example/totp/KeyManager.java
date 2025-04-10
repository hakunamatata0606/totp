package com.example.totp;

import android.content.SharedPreferences;
import android.content.Context;

public class KeyManager {
    private final SharedPreferences sharedPreferences;
    private static final String KeySharedPreferences = "this is key";

    public KeyManager(Context context) {
        sharedPreferences = context.getSharedPreferences(KeySharedPreferences, Context.MODE_PRIVATE);
    }

    public void set(String key, String value) {
        SharedPreferences.Editor editor = sharedPreferences.edit();
        editor.putString(key, value);
        editor.apply();
    }

    public String get(String key) {
        return sharedPreferences.getString(key, null);
    }
}
