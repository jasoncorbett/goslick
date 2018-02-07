import React, { Component } from 'react';

import App from 'grommet/components/App';
import Title from 'grommet/components/Title';
import Split from 'grommet/components/Split';
import Box from 'grommet/components/Box';
import Headline from 'grommet/components/Headline';
import Button from 'grommet/components/Button';
import PlatformGoogleIcon from 'grommet/components/icons/base/PlatformGoogle';
import TextInput from 'grommet/components/TextInput';
import FormField from 'grommet/components/FormField';
import Form from 'grommet/components/Form';
import Label from 'grommet/components/Label';
import FormFields from 'grommet/components/FormFields';

function isLoggedIn() {
    return localStorage.token;
}

export class SlickLogo extends Component {
  render() {
    return (
      <div>
        <Headline size={"xlarge"} style={{"fontFamily": "Audiowide, cursive"}}>S
          <span style={{"fontSize": "70%"}}>LICK</span>
        </Headline>
        <Headline size={"small"}>A Test Manager that doesn't suck!</Headline>
      </div>
    );
  }
}

export class LoginPage extends Component {
  render() {
    return (
      <Split flex="left">
        <Box justify="center" className="LoginPage" align="center"><SlickLogo/></Box>
        <Box justify="center" pad={{"horizontal": "small"}} full="vertical" align="center">
          <Button icon={<PlatformGoogleIcon/>} primary={true} href="/login/google" label="Login with Google"/>
        </Box>
      </Split>
    );
  }
}

export class UserInfoPage extends Component {
  constructor(props) {
    super(props);
    this.state = {
      value: "",
      message: ""
    };
    this.suggestions = ["hello", "world", "foo"];
  }

  handleClick(e) {
    e.preventDefault();
    delete localStorage.token;
    delete localStorage.userName;
    delete localStorage.userFirstName;
    delete localStorage.userFamilyName;
    delete localStorage.userGender;
    delete localStorage.userPicture;
    window.location.reload();
  }

  handleChange(event) {
    this.setState({value: event.target.value});
  }

  handleSubmit(event) {
    event.preventDefault();
    this.handlePermissionCheck();
  }

  handlePermissionCheck() {
    let that = this;
    fetch("/api/v1/isAuthorized/" + this.state.value, {headers: {"Authorization": "Bearer " + localStorage.token}}).then(function(response) {
      response.json().then(function(retval) {
        that.setState({"message": "You have permission " + that.state.value + ": " + retval.allowed});
      });
    });
  }

  handleSelect(selection) {
    this.setState({value: selection.suggestion});
  }

  render() {
    return (
      <Box>
        <Headline size="large">Welcome {localStorage.userFirstName}</Headline>
        <Button onClick={this.handleClick} label="Logout" />
        <Form onSubmit={this.handleSubmit.bind(this)} pad="large">
          <FormFields>
            <FormField>
              <TextInput onDOMChange={this.handleChange.bind(this)}
                         value={this.state.value}
                         onSelect={this.handleSelect.bind(this)}
                         suggestions={this.suggestions}/>
            </FormField>
          </FormFields>
          <Button type="submit" label="Check Permission" />
        </Form>
        <Label margin="large">{this.state.message}</Label>
      </Box>
    )
  }
}

export default class BasicApp extends Component {
  render() {
    let page = <LoginPage/>;
    if(isLoggedIn()) {
      page = <UserInfoPage/>;
    }
    return (
      <App centered={false}>
        {page}
      </App>
    );
  }
}
