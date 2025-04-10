import 'dart:async';
import 'dart:math';

import 'package:totp/config.dart';


// import 'package:flutter_secure_storage/flutter_secure_storage.dart';

// final storage = FlutterSecureStorage();

enum OtpStatusCode {
  ok(0),
  loginRequired(1),
  internalError(2),
  serverError(3),
  unknow(4);

  final int code;

  const OtpStatusCode(this.code);

  static OtpStatusCode fromInt(int code) {
    if (code < 0 || code >= OtpStatusCode.values.length) {
      throw ArgumentError("OtpStatusCode: invalid value");
    }
    return OtpStatusCode.values[code];
  }
}

class OtpResult {
  OtpStatusCode otpStatusCode;
  int otp;

  OtpResult({required this.otpStatusCode, required this.otp});

  static OtpResult fromMap(Map<Object?, Object?> map) {
    // todo: safely cast, indexing ?
    int statusCode = map['result'] as int;
    return OtpResult(otpStatusCode: OtpStatusCode.fromInt(statusCode), otp: map['otp'] as int);
  }
}


abstract class Totp {
  void Function(int) updateOtpCallback;
  void Function(OtpStatusCode) updateOtpStatusCallback;

  int _otpValue = 0;
  bool _initialized = false;

  late Timer _timer;

  Totp({required this.updateOtpCallback, required this.updateOtpStatusCallback});

  Future<OtpResult> getOtp();

  Future<OtpStatusCode> updateStatus(); 

  void start(int interval) {
    //Todo: handle time clock
    if (!_initialized) {
      _initialized = true;
    }
    _timer = Timer.periodic(Duration(seconds: interval), (timer) async {
      
      final OtpResult otpResult = await getOtp();
      // todo: handle result
      OtpStatusCode result = otpResult.otpStatusCode;
      if (result != OtpStatusCode.ok) {
        updateOtpStatusCallback(result);
        return;
      }
      int newOtp = otpResult.otp;

      if (newOtp != _otpValue) {
        _otpValue = newOtp;
        updateOtpCallback(_otpValue);
      }
    });
  }

  void stop() {
    if (_initialized) {
      _timer.cancel();
    }
  }
}

class TotpJava extends Totp {

  
  TotpJava(void Function(int) updateOtpCallback, void Function(OtpStatusCode) updateOtpStatusCallback)
    : super(updateOtpCallback: updateOtpCallback, updateOtpStatusCallback: updateOtpStatusCallback);

  @override
  Future<OtpResult> getOtp() async {
    final result = await totpPlatform.invokeMethod('generateOtp');
    if (result is Map) {
      final rs = OtpResult.fromMap(result);
      return rs;
    }
    return OtpResult(otpStatusCode: OtpStatusCode.internalError, otp: -1);
  }

  @override
  Future<OtpStatusCode> updateStatus() async {
    final result = await totpPlatform.invokeMethod('updateStatus');
    // todo: better implement ?
    updateOtpStatusCallback(OtpStatusCode.fromInt(result));
    return OtpStatusCode.fromInt(result);
  }
}