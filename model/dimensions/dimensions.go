/* Start of a package for managing SI units of dimension.

Copyright (C) Sarah Mount, 2011.

This program is free software; you can redistribute it and/or
modify it under the terms of the GNU General Public License
as published by the Free Software Foundation; either version 2
of the License, or (at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <http://www.gnu.org/licenses/>.
*/

package dimensions

// Standard units of dimension
// Map from description to unit.

/* 
 FIXME: This may be a better way to describe SI units. 
 FIXME: Should refactor, add symbolic manipulation and publish as separate package.

type si_unit string
const (
    Hz = "Hz"     
    radian = "radian" 
	steradian = "steradian"
    newton = "newton" 
    pascal = "pascal" 
    joule = "joule"  
    watt = "watt"   
    C = "C"      
    volt = "volt"   
    farad = "farad"  
    ohm = "ohm"    
    siemens = "siemens"
    webber = "webber" 
    tesla = "tesla"  
    henry = "henry"  
    Celsius = "Celsius"
    lumen = "lumen"  
    lux = "lux"    
	becquerel = "becquerel"
    gray = "gray"   
    sievert = "sievert"
    katal = "katal"  
	"Area"								: "m^2",
	"Volume"							: "m^3",
	"Velocity"							: "m/s",
	"Speed"								: "m/s",
	"Volumetric flow"					: "m^3/s",
	"Acceleration"						: "m/s^2",
	"Jerk"								: "m/s^3",
	"Jolt"								: "m/s^3",
	"Snap"								: "m/s^4",
    rad/s = "rad/s"  
    Ns = "Ns"     
    Nms = "Nms"    
    Nm = "Nm"     
    Kg = "Kg"     
    g = "g"      
 "J/m^3",  
 "kg/m^3", 
 "m^3/kg", 
 "mol/m^3",
 "m^3/mol",
 "J/K",    
 "J/kg",   
 "J/(Kmol)"
 "J/(Kkg)",
 "J/mol",  
 "W/m^2",  
 "W/m^2",  
 "W/(mk)", 
 "m^2/s",  
    Pas = "Pas"    
 "C/m^2",  
 "C/m^3",  
 "A/m^2",  
 "S/m",    
 "Sm^2/mol"
 "F/m",    
 "H/m",    
 "V/m",    
 "A/m",    
 "cd/m^3", 
 "C/kg",   
 "Gy/s",   
 "J/m^2",  
 "ohm/m",  
 "km/h",   
 "km/h",   
    percent = "percent"
	)
*/


var	si_units = map[string] string {
	"Hertz"								: "Hz",
	"Radian"							: "radian",
	"Steradian"							: "steradian",
	"Newton"							: "newton",
	"Pascal"							: "pascal",
	"Joule"								: "joule",
	"Watt"								: "watt",
	"Coulomb"							: "C",
	"Volt"								: "volt",
	"Farad"								: "farad",
	"Ohm"								: "ohm",
	"Siemen"							: "siemens",
	"Webber"							: "webber",
	"Tesla"								: "tesla",
	"Henry"								: "henry",
	"Celsius"							: "Celsius",
	"Lumen"								: "lumen",
	"Lux"								: "lux",
	"Becquerel"							: "becquerel",
	"Gray"								: "gray",
	"Sievert"							: "sievert",
	"Katal"								: "katal",
	"Area"								: "m^2",
	"Volume"							: "m^3",
	"Velocity"							: "m/s",
	"Speed"								: "m/s",
	"Volumetric flow"					: "m^3/s",
	"Acceleration"						: "m/s^2",
	"Jerk"								: "m/s^3",
	"Jolt"								: "m/s^3",
	"Snap"								: "m/s^4",
	"Angular velocity"					: "rad/s",
	"Momentum"							: "Ns",
	"Angular momentum"					: "Nms",
	"Newton meter"						: "Nm",
	"Kilogram"							: "Kg",
	"Gram"								: "g",
	"Energy density"					: "J/m^3",
	"Density"							: "kg/m^3",
	"Specific volume"					: "m^3/kg",
	"Concentration"						: "mol/m^3",
	"Molar volume"						: "m^3/mol",
	"Heat capacity"						: "J/K",
	"Specific energy"					: "J/kg",
	"Molar heat capacity"				: "J/(Kmol)",
	"Specific heat capacity"			: "J/(Kkg)",
	"Molar energy"						: "J/mol",
	"Heat Flux"							: "W/m^2",
	"Irradiance"						: "W/m^2",
	"Thermal conductivity"				: "W/(mk)",
	"Kinematic viscosity"				: "m^2/s",
	"Dynamic viscosity"					: "Pas",
	"Electric displacement field"		: "C/m^2",
	"Electric charge density"			: "C/m^3",
	"Electric Current density"			: "A/m^2",
	"Conductivity"						: "S/m",
	"Molar conductivity"				: "Sm^2/mol",
	"Permittivity"						: "F/m",
	"Permeability"						: "H/m",
	"Electric field"					: "V/m",
	"Magnetic field"					: "A/m",
	"Luminance"							: "cd/m^3",
	"Exposure"							: "C/kg",
	"Absorbed dose"						: "Gy/s",
	"Surface tension"					: "J/m^2",
	"Resistivity"						: "ohm/m",
//	"Speed"								: "km/h",
//	"Velocity"							: "km/h",
	"Percentage"						: "percent",
}