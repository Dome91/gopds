import {ComponentFixture, TestBed} from '@angular/core/testing';

import {SettingsUsersComponent} from './settings-users.component';
import {ToastrModule} from "ngx-toastr";
import {By} from "@angular/platform-browser";
import {of} from "rxjs";
import {FontAwesomeModule} from "@fortawesome/angular-fontawesome";
import createSpyObj = jasmine.createSpyObj;
import {UserService} from "../../../services/user.service";
import SpyObj = jasmine.SpyObj;
import {Role, User} from "../../../models/user";

describe('SettingsUsersComponent', () => {
  let component: SettingsUsersComponent;
  let fixture: ComponentFixture<SettingsUsersComponent>;

  let userService: SpyObj<UserService>;
  let users = [new User('user1', 'password1', Role.ADMIN), new User('user2', 'password2', Role.USER)]

  beforeEach(async () => {
    userService = createSpyObj<UserService>(['getAll', 'delete']);
    await TestBed.configureTestingModule({
      declarations: [SettingsUsersComponent],
      imports: [ToastrModule.forRoot(), FontAwesomeModule],
      providers: [{provide: UserService, useValue: userService}]
    })
      .compileComponents();
  });

  beforeEach(() => {
    userService.getAll.and.returnValue(of(users));
    fixture = TestBed.createComponent(SettingsUsersComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });

  it('should render users', () => {
    expect(fixture.debugElement.query(By.css('#username-user1'))).toBeTruthy();
    expect(fixture.debugElement.query(By.css('#username-user2'))).toBeTruthy();
  });

  it('should delete user', () => {
    userService.delete.and.returnValue(of(null));
    const createButton = fixture.debugElement.query(By.css('#delete-user1'));
    createButton.triggerEventHandler('click', null);
    expect(userService.delete).toHaveBeenCalledWith(users[0]);
  });
});
