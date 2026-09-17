import { HttpClient } from '@angular/common/http';
import { inject, Service } from '@angular/core';
import { Observable } from 'rxjs';
import { InvitationTokenRequest, InvitationTokenResponse } from '../models/invite.model';

@Service()
export class InvitationService {
    private http = inject(HttpClient);

    public getInvitationToken(payload: InvitationTokenRequest): Observable<InvitationTokenResponse> {

        return this.http.post<InvitationTokenResponse>(
            '/api/auth/invite-token/create',
            payload
        )
    }
}
