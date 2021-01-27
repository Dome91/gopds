import {Component, OnInit} from '@angular/core';
import {Role, User} from "../../../../models/user";
import {UserService} from "../../../../services/user.service";
import {Router} from "@angular/router";

@Component({
  selector: 'app-create-user',
  templateUrl: './create-user.component.html',
  styleUrls: ['./create-user.component.sass']
})
export class CreateUserComponent implements OnInit {

  user: User;
  Role = Role;

  constructor(private userService: UserService, private router: Router) {
    this.user = new User('', '', Role.USER);
  }

  ngOnInit(): void {
  }

  create() {
    this.userService.create(this.user).subscribe(
      () => this.router.navigateByUrl('/settings'),
    );
  }
}
