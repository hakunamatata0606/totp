import 'dart:async';

import 'package:flutter/material.dart';
import 'package:totp/login.dart';
import 'totp.dart';
import 'package:intl/intl.dart';

void main() {
  runApp(const MainApp());
}

class MainApp extends StatelessWidget {
  const MainApp({super.key});

  @override
  Widget build(BuildContext context) {
    return const MaterialApp(
      home: MyHome(),
    );
  }
}

enum TotpState {
  initial,
  login,
  idle
}

class MyHome extends StatefulWidget {
  const MyHome({Key? key}) : super(key: key);
  
  @override
  MyHomeState createState() {
      return MyHomeState();
  }
}

class MyHomeState extends State<MyHome> {

  String _otpValue = "";
  TotpState _state = TotpState.initial;

  late Totp _totp;

  
  @override
  void initState() {
    super.initState();
    _startTotp();
  }

  void _startTotp() {
    _totp = TotpJava (
      onUpdateOtpValue,
      onUpdateStatus
    );
    _totp.updateStatus();
  }

  void onUpdateOtpValue(int otp) {
    setState(() {
      var formatOtp = NumberFormat('000,000').format(otp);
      _otpValue = formatOtp;
    });
  }

  void onUpdateStatus(OtpStatusCode status) {  
    if (status == OtpStatusCode.loginRequired) {
      if (_state == TotpState.idle) {
        _totp.stop();
      }
      updateState(TotpState.login);
    }else if (status == OtpStatusCode.ok) {
      _totp.stop();
      //todo: get from config file ?
      _totp.start(5);
      updateState(TotpState.idle);
    }
  }

  void onLoggedIn() {
    _totp.stop();
    _totp.start(5);
    updateState(TotpState.idle);
  }

  void updateState(TotpState newState) {
    setState(() {
      _state = newState;
    });
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        backgroundColor: Colors.orange,
        title: const Text("TOTP"),
      ),
      body: Center(
        child: Builder(
          builder: (context) {
            switch (_state) {
              case TotpState.initial:
                return const CircularProgressIndicator();
              case TotpState.login:
                return Login(loggedInCallback: onLoggedIn);
              case TotpState.idle:
                return Container(
                  alignment: Alignment.center,
                  padding: const EdgeInsets.all(8.0),
                  child:  Column(
                    children: [
                      const Text(
                        "Generate otp",
                        style: TextStyle(
                          fontSize: 32,
                          fontWeight: FontWeight.bold
                        ),
                      ),
                      Text(
                        _otpValue,
                        style: const TextStyle(
                          color: Colors.blue,
                          fontSize: 48,
                          fontWeight: FontWeight.bold
                        ),
                      )
                    ],
                  ),
                );
            }
          }
          ),
        )
      
    );
  }
}