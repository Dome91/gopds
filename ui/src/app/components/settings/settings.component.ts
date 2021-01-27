import { Component, OnInit } from '@angular/core';
import {faBookOpen, faPlus, faUsers} from "@fortawesome/free-solid-svg-icons";

@Component({
  selector: 'app-settings',
  templateUrl: './settings.component.html',
  styleUrls: ['./settings.component.sass']
})
export class SettingsComponent implements OnInit {

  faBookOpen = faBookOpen;
  faPlus = faPlus;
  faUsers = faUsers;

  constructor() { }

  ngOnInit(): void {
  }

}
