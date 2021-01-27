import {ComponentFixture, TestBed, waitForAsync} from '@angular/core/testing';

import {CreateSourceComponent} from './create-source.component';
import {SourceService} from "../../../../services/source.service";
import SpyObj = jasmine.SpyObj;
import {Router} from "@angular/router";
import createSpyObj = jasmine.createSpyObj;
import {MockNavbarComponent} from "../../../mocks/mock-navbar.components";
import {ToastrModule} from "ngx-toastr";
import {FormsModule} from "@angular/forms";
import {of} from "rxjs";
import {By} from "@angular/platform-browser";
import {Source} from "../../../../models/source";

describe('CreateSourceComponent', () => {
  let component: CreateSourceComponent;
  let fixture: ComponentFixture<CreateSourceComponent>;

  let sourceService: SpyObj<SourceService>;
  let router: SpyObj<Router>;

  beforeEach(async () => {
    sourceService = createSpyObj<SourceService>(['create']);
    router = createSpyObj<Router>(['navigateByUrl']);

    await TestBed.configureTestingModule({
      declarations: [CreateSourceComponent, MockNavbarComponent],
      imports: [FormsModule, ToastrModule.forRoot()],
      providers: [
        {provide: SourceService, useValue: sourceService},
        {provide: Router, useValue: router}
      ]
    })
      .compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(CreateSourceComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });

  it('should create a source', waitForAsync(() => {
    sourceService.create.and.returnValue(of(null));

    fixture.whenStable().then(() => {
      const nameInput = fixture.debugElement.query(By.css('#name')).nativeElement;
      const pathInput = fixture.debugElement.query(By.css('#path')).nativeElement;

      nameInput.value = 'source1'
      nameInput.dispatchEvent(new Event('input'))

      pathInput.value = '/source1/path'
      pathInput.dispatchEvent(new Event('input'))

      const createButton = fixture.debugElement.query(By.css('#createButton'));
      createButton.triggerEventHandler('click', null);

      expect(sourceService.create).toHaveBeenCalledWith(new Source('', 'source1', '/source1/path'));
      expect(router.navigateByUrl).toHaveBeenCalledWith('/settings');
    });
  }));
});
