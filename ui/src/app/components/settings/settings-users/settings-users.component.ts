import {Component, OnInit} from '@angular/core';
import {faEdit, faTrash, faUser, faUserShield} from "@fortawesome/free-solid-svg-icons";
import {Role, User} from "../../../models/user";
import {UserService} from "../../../services/user.service";

@Component({
  selector: 'app-settings-users',
  templateUrl: './settings-users.component.html',
  styleUrls: ['./settings-users.component.sass']
})
export class SettingsUsersComponent implements OnInit {

  users: User[];

  faTrash = faTrash;
  faEdit = faEdit
  faUser = faUser;
  faUserShield = faUserShield
  Role = Role;

  constructor(private userService: UserService) {
    this.users = [];
  }

  ngOnInit(): void {
    this.fetchAllUsers();
  }

  private fetchAllUsers() {
    this.userService.getAll().subscribe(
      users => this.users = users
    );
  }

  delete(user: User) {
    this.userService.delete(user).subscribe(
      () => this.fetchAllUsers()
    );
  }

  isAdmin(user: User): boolean {
    return user.role === Role.ADMIN
  }

  isUser(user: User): boolean {
    return user.role === Role.USER
  }

}
