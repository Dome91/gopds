import { Component, OnInit } from '@angular/core';
import {HttpErrorResponse} from "@angular/common/http";
import {Router} from "@angular/router";
import {LoginService} from "../../services/login.service";
import {ToastrService} from "ngx-toastr";

@Component({
  selector: 'app-login',
  templateUrl: './login.component.html',
  styleUrls: ['./login.component.sass']
})
export class LoginComponent implements OnInit {

  username: string;
  password: string;

  constructor(private loginService: LoginService,
              private toastrService: ToastrService,
              private router: Router) {
    this.username = "";
    this.password = "";
  }

  ngOnInit(): void {
    //this.isLoggedIn();
  }

  isLoggedIn() {
    this.loginService.isLoggedIn().subscribe(
      () => this.router.navigateByUrl('/catalog'),
      () => {
      }
    )
  }

  login(): void {
    this.loginService.login(this.username, this.password)
      .subscribe(() => this.router.navigateByUrl('/catalog'),
        (_: HttpErrorResponse) => this.toastrService.error('Login failed.'));
  }

}
