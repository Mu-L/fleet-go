module.exports = {


  friendlyName: 'Remove one compliance partner tenant',


  description: 'Updates a microsfot compliance tenant\'s status as "deprovisioned" and deletes the associated Database record',

  inputs: {
    entraTenantId: {
      type: 'string',
      required: true,
    },
    fleetServerSecret: {
      type: 'string',
      requried: true,
    },
  },


  exits: {
    success: {
      description: 'The requesting entra tenant has been successfully deprovisioned.'
    }
  },


  fn: async function ({entraTenantId, fleetServerSecret}) {

    let informationAboutThisTenant = await MicrosoftComplianceTenant.findOne({entraTenantId: entraTenantId, fleetServerSecret: fleetServerSecret});
    if(!informationAboutThisTenant) {
      return new Error({error: 'No MicrosoftComplianceTenant record was found that matches the provided entra_tenant_id and fleet_server_secret combination.'});// TODO: return a more clear error.
    }

    let tokenAndApiUrls = await sails.helpers.microsoftProxy.getAccessTokenAndApiUrls.with({
      complianceTenantRecordId: informationAboutThisTenant.id
    });

    let accessToken = tokenAndApiUrls.accessToken;
    let tenantDataSyncUrl = tokenAndApiUrls.tenantDataSyncUrl;


    // Deprovison this tenant
    await sails.helpers.http.sendHtttpRequest.with({
      method: 'PUT',
      url: `${tenantDataSyncUrl}/${encodeURIComponent(`PartnerTenants(guid${informationAboutThisTenant.entraTenantId}`)}?api-version=1.0`,
      headers: {
        'Authorization': `Bearer ${accessToken}`
      },
      body: {
        Provisioned: 2,// 1 = provisioned, 2 = deprovisioned.
        PartnerEnrollmentUrl: '', //TODO: how do we get this, the example in microsoft's docs are using customer.com/enrollment, so does this need to be a value of a url on the connected Fleet instance?
        PartnerRemediationUrl: '', // TODO: same as the above.
      }
    }).intercept((err)=>{
      return new Error({error: `an error occurred when deprovisioning a Microsoft compliance tenant. Full error: ${require('util').inspect(err, {depth: 3})}`});
    });

    await MicrosoftComplianceTenant.destroyOne({id: informationAboutThisTenant.id});


    // All done.
    return;

  }


};
