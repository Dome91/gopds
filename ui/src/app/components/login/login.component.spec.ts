import {ComponentFixture, TestBed, waitForAsync} from '@angular/core/testing';

import { LoginComponent } from './login.component';
import {of, throwError} from "rxjs";
import {FormsModule} from "@angular/forms";
import {ToastrModule, ToastrService} from "ngx-toastr";
import {LoginService} from "../../services/login.service";
import {Router} from "@angular/router";
import {BrowserAnimationsModule} from "@angular/platform-browser/animations";
import {By} from "@angular/platform-browser";
import SpyObj = jasmine.SpyObj;

describe('LoginComponent', () => {
  let component: LoginComponent;
  let fixture: ComponentFixture<LoginComponent>;

  let loginService: SpyObj<LoginService>;
  let router: SpyObj<Router>;

  function login() {
    const usernameInput = fixture.debugElement.query(By.css('#username')).nativeElement;
    const passwordInput = fixture.debugElement.query(By.css('#password')).nativeElement;

    usernameInput.value = 'callisto';
    usernameInput.dispatchEvent(new Event('input'));

    passwordInput.value = 'callistoPassword';
    passwordInput.dispatchEvent(new Event('input'));

    const loginButton = fixture.debugElement.query(By.css('#loginButton'));
    loginButton.nativeElement.click();
  }

  beforeEach(async () => {
    loginService = jasmine.createSpyObj<LoginService>(['login', 'isLoggedIn']);
    router = jasmine.createSpyObj<Router>(['navigateByUrl']);

    await TestBed.configureTestingModule({
      declarations: [LoginComponent],
      imports: [FormsModule, BrowserAnimationsModule, ToastrModule.forRoot()],
      providers: [ToastrService,
        {provide: LoginService, useValue: loginService},
        {provide: Router, useValue: router}]
    }).compileComponents();
  });

  beforeEach(() => {
    loginService.isLoggedIn.and.returnValue(throwError('Not logged in.'))
    fixture = TestBed.createComponent(LoginComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });

  it('should redirect for successful login', waitForAsync(() => {
    loginService.login.and.returnValue(of(null));

    fixture.whenStable().then(() => {
      login();
      expect(router.navigateByUrl).toHaveBeenCalledWith('/catalog');
      expect(loginService.login).toHaveBeenCalledWith('callisto', 'callistoPassword');
    });
  }));

  it('should not redirect when login fails', waitForAsync(() => {
    loginService.login.and.returnValue(throwError('Login failed'));

    fixture.whenStable().then(() => {
      login();
      expect(router.navigateByUrl).not.toHaveBeenCalled();
      expect(loginService.login).toHaveBeenCalledWith('callisto', 'callistoPassword');
    });
  }));

  /*it('should redirect for user already logged in', () => {
    loginService.isLoggedIn.and.returnValue(of(null));
    component.ngOnInit();
    expect(router.navigateByUrl).toHaveBeenCalledWith('/catalog');
  });*/
});
