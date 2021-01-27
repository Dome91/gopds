import { ComponentFixture, TestBed } from '@angular/core/testing';

import { NavbarComponent } from './navbar.component';
import {By} from "@angular/platform-browser";
import {of} from "rxjs";
import {LoginService} from "../../services/login.service";
import {Router} from "@angular/router";
import {ToastrModule} from "ngx-toastr";
import createSpyObj = jasmine.createSpyObj;
import SpyObj = jasmine.SpyObj;

describe('NavbarComponent', () => {
  let component: NavbarComponent;
  let fixture: ComponentFixture<NavbarComponent>;

  let router: SpyObj<Router>;
  let loginService: SpyObj<LoginService>;

  beforeEach(async () => {
    router = createSpyObj<Router>(['navigateByUrl']);
    loginService = createSpyObj<LoginService>(['logout']);

    await TestBed.configureTestingModule({
      declarations: [NavbarComponent],
      imports: [ToastrModule.forRoot()],
      providers: [{provide: Router, useValue: router},
        {provide: LoginService, useValue: loginService}]
    })
      .compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(NavbarComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });

  it('should route correctly', () => {
    const brandLink = fixture.debugElement.query(By.css('#brand-link')).nativeElement;
    const seriesLink = fixture.debugElement.query(By.css('#catalog-link')).nativeElement;
    const settingsLink = fixture.debugElement.query(By.css('#settings-link')).nativeElement;
    expect(brandLink.getAttribute('routerLink')).toEqual('/catalog')
    expect(seriesLink.getAttribute('routerLink')).toEqual('/catalog')
    expect(settingsLink.getAttribute('routerLink')).toEqual('/settings')
  });

  it('should logout', () => {
    loginService.logout.and.returnValue(of(null));
    component.logout();
    expect(loginService.logout).toHaveBeenCalled();
    expect(router.navigateByUrl).toHaveBeenCalledWith('/');
  });
});
