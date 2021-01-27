import {ComponentFixture, TestBed, waitForAsync} from '@angular/core/testing';

import {CreateUserComponent} from './create-user.component';
import {Router} from "@angular/router";
import {MockNavbarComponent} from "../../../mocks/mock-navbar.components";
import {FormsModule} from "@angular/forms";
import {of} from "rxjs";
import {By} from "@angular/platform-browser";
import {UserService} from "../../../../services/user.service";
import {Role, User} from "../../../../models/user";
import SpyObj = jasmine.SpyObj;
import createSpyObj = jasmine.createSpyObj;

describe('CreateUserComponent', () => {
  let component: CreateUserComponent;
  let fixture: ComponentFixture<CreateUserComponent>;

  let userService: SpyObj<UserService>;
  let router: SpyObj<Router>;

  beforeEach(async () => {
    userService = createSpyObj<UserService>(['create']);
    router = createSpyObj<Router>(['navigateByUrl']);

    await TestBed.configureTestingModule({
      declarations: [CreateUserComponent, MockNavbarComponent],
      imports: [FormsModule],
      providers: [
        {provide: UserService, useValue: userService},
        {provide: Router, useValue: router}
      ]
    })
      .compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(CreateUserComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });

  it('should create a user', waitForAsync(() => {
    userService.create.and.returnValue(of(null));

    fixture.whenStable().then(() => {
      createUser(() => {
      }, Role.USER);
    })
  }));

  it('should create an admin', waitForAsync(() => {
    userService.create.and.returnValue(of(null));

    fixture.whenStable().then(() => {

      createUser(() => {
        fixture.debugElement.query(By.css('#dropdownAdminButton'))
          .triggerEventHandler('click', null);
      }, Role.ADMIN)
    })
  }));

  function createUser(setRoleFn: Function, expectedRole: Role) {
    const usernameInput = fixture.debugElement.query(By.css('#username')).nativeElement;
    const passwordInput = fixture.debugElement.query(By.css('#password')).nativeElement;
    setRoleFn();

    usernameInput.value = 'user1'
    usernameInput.dispatchEvent(new Event('input'))

    passwordInput.value = 'password1'
    passwordInput.dispatchEvent(new Event('input'))

    const createButton = fixture.debugElement.query(By.css('#createButton'));
    createButton.triggerEventHandler('click', null);

    expect(userService.create).toHaveBeenCalledWith(new User('user1', 'password1', expectedRole));
    expect(router.navigateByUrl).toHaveBeenCalledWith('/settings');
  }
});
