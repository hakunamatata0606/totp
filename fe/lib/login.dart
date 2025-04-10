import 'package:flutter/material.dart';
import 'package:totp/config.dart';
import 'package:totp/totp.dart';

class Login extends StatefulWidget {
  final void Function() loggedInCallback;

  const Login({Key? key, required this.loggedInCallback}) : super(key: key);

  @override
  LoginState createState() => LoginState();
}

class LoginState extends State<Login> {
  late void Function() loggedInCallback;
  final TextEditingController _usernameController = TextEditingController();
  final TextEditingController _passwordController = TextEditingController();
  String _additionalMessage = "";

  @override
  void initState() {
    super.initState();
    loggedInCallback = widget.loggedInCallback;
  }

  void _login() async {
    String username = _usernameController.text;
    String password = _passwordController.text;

    final result = await totpPlatform.invokeMethod("login", {'username': username, 'password': password});
    final resultCode = OtpStatusCode.fromInt(result);
    if (resultCode == OtpStatusCode.ok) {
      loggedInCallback();
    }else if (resultCode == OtpStatusCode.loginRequired) {
      setMessage("Username or password is not correct");
    }else {
      setMessage("Internal error");
    }
  }

  void setMessage(String message) {
    setState(() {
      _additionalMessage = message;
    });
  }

  @override
  Widget build(BuildContext context) {
    return Padding(
      padding: const EdgeInsets.all(16),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          const Text('Username:'),
          const SizedBox(height: 10),
          TextField(
            controller: _usernameController,
            decoration: const InputDecoration(
              labelText: "Enter username",
              border: OutlineInputBorder()
            ),
            maxLines: 1,
          ),
          const SizedBox(height: 20),
          const Text('Password:'),
          const SizedBox(height: 10),
          TextField(
            obscureText: true,
            controller: _passwordController,
            decoration: const InputDecoration(
              labelText: "Enter password",
              border: OutlineInputBorder()
            ),
            maxLines: 1,
          ),
          const SizedBox(height: 20),
          if (_additionalMessage.isNotEmpty)
            Container(
              alignment: Alignment.center,
              child: Column(children: [
                Text(_additionalMessage, style: const TextStyle(color: Colors.red),),
                const SizedBox(height: 20),
              ]),
            ),
          ElevatedButton(
            onPressed: () {
              _login();
            },
            style: ElevatedButton.styleFrom(
              minimumSize: const Size(double.infinity, 50)
            ),
            child: const Text('Login')
          )
        ],
      ), 
    );
  }
}