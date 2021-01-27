import { Component, OnInit } from '@angular/core';
import {Router} from "@angular/router";
import {LoginService} from "../../services/login.service";

@Component({
  selector: 'app-navbar',
  templateUrl: './navbar.component.html',
  styleUrls: ['./navbar.component.sass']
})
export class NavbarComponent implements OnInit {

  constructor(private router: Router, private loginService: LoginService) {
  }

  ngOnInit(): void {
  }

  logout() {
    this.loginService.logout().subscribe(
      () => this.router.navigateByUrl('/'),
    );
  }

}
