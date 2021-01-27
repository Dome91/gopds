import { ComponentFixture, TestBed } from '@angular/core/testing';

import { SettingsComponent } from './settings.component';
import {FontAwesomeModule} from "@fortawesome/angular-fontawesome";
import {Component} from "@angular/core";
import {MockNavbarComponent} from "../mocks/mock-navbar.components";

describe('SettingsComponent', () => {
  let component: SettingsComponent;
  let fixture: ComponentFixture<SettingsComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [SettingsComponent, MockSettingsLibraryComponent, MockNavbarComponent],
      imports: [FontAwesomeModule]
    })
      .compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(SettingsComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});

@Component({
  selector: 'app-settings-library',
  template: '<p>Mock Settings Library Component</p>'
})
class MockSettingsLibraryComponent {
}
